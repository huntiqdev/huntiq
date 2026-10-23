package fpa

import (
	"math"
)

// WorkloadMetrics calculates demand from user activity
type WorkloadMetrics struct {
	DailyRequests        float64
	AverageRPS           float64
	PeakRPS              float64
	MonthlyRequests      float64
	DBReadIOPS           float64
	DBWriteIOPS          float64
	AccumulatedStorageTB float64
	MonthlyBandwidthTB   float64
	DailyBackgroundJobs  float64
}

// CalculateWorkload translates active user metrics into technical infrastructure demand
func CalculateWorkload(dau, mau, cumulativeUsers float64, monthIndex int, cfg WorkloadConfig) WorkloadMetrics {
	if dau < 0 {
		dau = 0
	}
	sessionsPerUser := cfg.SessionsPerUserPerDay
	if sessionsPerUser <= 0 {
		sessionsPerUser = 2.5
	}
	reqsPerSession := cfg.RequestsPerSession
	if reqsPerSession <= 0 {
		reqsPerSession = 18.0
	}
	peakRatio := cfg.PeakToAverageRatio
	if peakRatio <= 1.0 {
		peakRatio = 3.0
	}

	dailyRequests := dau * sessionsPerUser * reqsPerSession
	avgRPS := dailyRequests / 86400.0
	peakRPS := avgRPS * peakRatio
	monthlyRequests := dailyRequests * 30.42

	dbReadsPerReq := cfg.DBReadsPerRequest
	if dbReadsPerReq <= 0 {
		dbReadsPerReq = 3.5
	}
	dbWritesPerReq := cfg.DBWritesPerRequest
	if dbWritesPerReq <= 0 {
		dbWritesPerReq = 0.4
	}

	dbReadIOPS := avgRPS * dbReadsPerReq
	dbWriteIOPS := avgRPS * dbWritesPerReq

	// Image and Object storage accumulation
	avgImgKB := cfg.AvgImageSizeKB
	if avgImgKB <= 0 {
		avgImgKB = 350.0 // 350 KB compressed WebP
	}
	uploadsPerUser := cfg.ImageUploadsPerActiveUser
	if uploadsPerUser <= 0 {
		uploadsPerUser = 0.05
	}

	// Monthly new image storage in Gigabytes
	monthlyImageUploads := mau * uploadsPerUser * 4.0 // 4 photos per listing
	monthlyImageGB := (monthlyImageUploads * avgImgKB) / (1024.0 * 1024.0)

	// Cumulative storage over months + 3x replication & backup
	accumulatedStorageGB := (monthlyImageGB * float64(monthIndex)) * 1.5
	// Database structured data: users + products + transactions + logs
	dbDataGB := (cumulativeUsers * (cfg.DataRetainedPerUserKB / (1024.0 * 1024.0))) + (float64(monthIndex) * 2.0)
	totalStorageTB := (accumulatedStorageGB + dbDataGB) / 1024.0

	// Bandwidth egress: listing page views + image views + API JSON payloads
	imgViewsPerUser := cfg.ImageViewsPerUserDay
	if imgViewsPerUser <= 0 {
		imgViewsPerUser = 24.0
	}
	dailyBandwidthGB := (dau * imgViewsPerUser * (avgImgKB * 0.7) / (1024.0 * 1024.0)) + // 70% cache hit or thumbnail
		(dailyRequests * 2.5 / (1024.0 * 1024.0)) // 2.5 KB avg API payload
	monthlyBandwidthTB := (dailyBandwidthGB * 30.42) / 1024.0

	jobsPerUser := cfg.BackgroundJobsPerUserDay
	if jobsPerUser <= 0 {
		jobsPerUser = 0.20
	}
	dailyBackgroundJobs := dau * jobsPerUser

	return WorkloadMetrics{
		DailyRequests:        dailyRequests,
		AverageRPS:           avgRPS,
		PeakRPS:              peakRPS,
		MonthlyRequests:      monthlyRequests,
		DBReadIOPS:           dbReadIOPS,
		DBWriteIOPS:          dbWriteIOPS,
		AccumulatedStorageTB: totalStorageTB,
		MonthlyBandwidthTB:   monthlyBandwidthTB,
		DailyBackgroundJobs:  dailyBackgroundJobs,
	}
}

// Assess20MCapacityScenario performs dedicated workload capacity stress analysis for 20M DAU
func Assess20MCapacityScenario(workloadCfg WorkloadConfig, infraCfg InfraCostConfig) CapacityPlanning20M {
	targetDAU := 20000000.0
	dauToMAU := workloadCfg.DAUtoMAURatio
	if dauToMAU <= 0 {
		dauToMAU = 0.25
	}
	estimatedMAU := targetDAU / dauToMAU // 80,000,000 MAU

	metrics := CalculateWorkload(targetDAU, estimatedMAU, 120000000.0, 60, workloadCfg)

	// App server compute sizing:
	// Go with Fiber can handle ~8,000-15,000 simple RPS per modern x86/ARM core under optimized conditions.
	// At Peak RPS, with 50% target utilization buffer:
	coresNeeded := int(math.Ceil((metrics.PeakRPS / 5000.0) * 1.5))
	ramGBNeeded := int(math.Ceil(float64(coresNeeded) * 4.0)) // 4GB RAM per core

	// Cloud cost estimate for 20M DAU scale:
	monthlyCloud := (float64(coresNeeded) * 1200.0) + // Compute vCPU & RAM in ETB
		(metrics.AccumulatedStorageTB * 1024.0 * infraCfg.ObjectStorageCostPerGB) +
		(metrics.MonthlyBandwidthTB * 1024.0 * infraCfg.CDNEgressCostPerGBETB) +
		(metrics.DBReadIOPS * 0.15) + 350000.0 // Managed DB cluster, Redis cluster, WAF

	return CapacityPlanning20M{
		TargetDAU:               targetDAU,
		EstimatedMAU:            estimatedMAU,
		PeakRPS:                 metrics.PeakRPS,
		AverageRPS:              metrics.AverageRPS,
		DailyAPIRequests:        metrics.DailyRequests,
		MonthlyAPIRequests:      metrics.MonthlyRequests,
		RequiredAppCores:        coresNeeded,
		RequiredRAMGB:           ramGBNeeded,
		DatabaseReadIOPS:        metrics.DBReadIOPS,
		DatabaseWriteIOPS:       metrics.DBWriteIOPS,
		EstimatedListingImagesTB: metrics.AccumulatedStorageTB,
		MonthlyEgressBandwidthTB: metrics.MonthlyBandwidthTB,
		EstimatedMonthlyCloudETB: monthlyCloud,
		ArchitectureFeasibility: "EXTREME DEMAND CAPACITY SCENARIO — Distributed Architecture Required",
		TechnicalBottlenecks: []string{
			"Single-instance PostgreSQL cannot sustain ~250k+ write IOPS and ~2M+ read IOPS without horizontal sharding or read replica clusters.",
			"Peak RPS exceeds 700k req/s; requires Cloudflare/Fastly CDN edge caching (min 85% cache hit ratio) to protect origin application servers.",
			"Object storage requires distributed MinIO cluster or multi-region S3 with tiering to cold storage (infrequent access).",
			"Centralized Redis session store must be migrated to a partitioned Redis Cluster with replica read offload.",
			"Direct database migrations become dangerous at 100M+ rows; requires zero-downtime shadow table schema migrations (e.g. pg_repack / gh-ost style).",
		},
		RequiredScaleMeasures: []string{
			"Implement edge caching for public product feeds, category hierarchies, and static assets.",
			"Deploy PostgreSQL read replicas for read-heavy search and catalog listing queries.",
			"Implement horizontal connection pooling (PgBouncer) with transaction-level pooling across all app workers.",
			"Offload image processing to asynchronous worker queues (RabbitMQ/Kafka) with dedicated GPU/CPU transcoders.",
			"Conduct rigorous end-to-end distributed load tests with tools like k6 / Locust before claiming production readiness.",
		},
	}
}
