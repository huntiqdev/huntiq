package fpa

import (
	"math"
)

// InfraCostBreakdown holds all itemized hosting costs for a month
type InfraCostBreakdown struct {
	ComputeAndNetwork float64
	Database          float64
	StorageAndCDN     float64
	SecurityAndAuth   float64
	LoggingAndBackup  float64
	TotalMonthlyCost  float64
}

// CalculateInfraCost computes monthly infrastructure bills from workload and pricing parameters
func CalculateInfraCost(workload WorkloadMetrics, activeUsers float64, monthIndex int, cfg InfraCostConfig) InfraCostBreakdown {
	// Profile scaling factor
	profileMultiplier := 1.0
	switch cfg.Profile {
	case ProfileLeanEarlyStage:
		profileMultiplier = 0.40
	case ProfileHighAvailability:
		profileMultiplier = 1.75
	case ProfileLargeMultiRegion:
		profileMultiplier = 3.20
	case ProfileGrowingProduction:
		fallthrough
	default:
		profileMultiplier = 1.0
	}

	// 1. Application Compute & Load Balancing
	baseCompute := cfg.ComputeBaseMonthlyETB
	if baseCompute <= 0 {
		baseCompute = 4500.0 // ETB
	}
	rpsCostRate := cfg.ComputeCostPer1KRPS
	if rpsCostRate <= 0 {
		rpsCostRate = 250.0 // ETB per 1k Peak RPS
	}
	peakKiloRPS := workload.PeakRPS / 1000.0
	computeCost := (baseCompute + (peakKiloRPS * rpsCostRate)) * profileMultiplier

	// 2. Database (PostgreSQL compute, replica, and disk)
	baseDB := cfg.DBBaseMonthlyETB
	if baseDB <= 0 {
		baseDB = 6000.0 // ETB
	}
	dbStorageRate := cfg.DBStorageCostPerGBETB
	if dbStorageRate <= 0 {
		dbStorageRate = 18.0 // ETB per GB SSD
	}
	dbDataGB := (workload.AccumulatedStorageTB * 1024.0) * 0.15 // DB structured slice
	dbCost := (baseDB + (dbDataGB * dbStorageRate) + (workload.DBWriteIOPS * 0.05)) * profileMultiplier

	// 3. Object Storage (Images) & CDN Egress
	objStorageRate := cfg.ObjectStorageCostPerGB
	if objStorageRate <= 0 {
		objStorageRate = 3.50 // ETB per GB per month (e.g. S3 / MinIO)
	}
	cdnRate := cfg.CDNEgressCostPerGBETB
	if cdnRate <= 0 {
		cdnRate = 8.50 // ETB per GB egress
	}
	totalStorageGB := workload.AccumulatedStorageTB * 1024.0
	storageCost := (totalStorageGB * objStorageRate)
	cdnCost := (workload.MonthlyBandwidthTB * 1024.0 * cdnRate)
	storageAndCDNCost := storageCost + cdnCost

	// 4. Security, WAF, SMS, Email, Auth
	wafBase := cfg.SecurityWAFMonthlyETB
	if wafBase <= 0 {
		wafBase = 3000.0 // ETB
	}
	smsUnitRate := cfg.AuthEmailSMSUnitETB
	if smsUnitRate <= 0 {
		smsUnitRate = 1.20 // ETB per OTP / notification SMS
	}
	monthlySMS := activeUsers * 0.40 // 40% of users trigger login OTP / order notification monthly
	securityAuthCost := (wafBase + (monthlySMS * smsUnitRate)) * profileMultiplier

	// 5. Logging, Monitoring & Backup
	logBase := cfg.LoggingMonitoringETB
	if logBase <= 0 {
		logBase = 2500.0
	}
	backupBase := cfg.BackupDisasterRecETB
	if backupBase <= 0 {
		backupBase = 2000.0
	}
	logBackupCost := (logBase + backupBase + (totalStorageGB * 0.20)) * profileMultiplier

	subtotal := computeCost + dbCost + storageAndCDNCost + securityAuthCost + logBackupCost

	// 6. Redundancy & Capacity Reserve Buffer
	bufferRate := cfg.RedundancyBufferPct
	if bufferRate <= 0 {
		bufferRate = 0.20 // 20% reserve headroom
	}
	totalWithBuffer := subtotal * (1.0 + bufferRate)

	return InfraCostBreakdown{
		ComputeAndNetwork: round2(computeCost),
		Database:          round2(dbCost),
		StorageAndCDN:     round2(storageAndCDNCost),
		SecurityAndAuth:   round2(securityAuthCost),
		LoggingAndBackup:  round2(logBackupCost),
		TotalMonthlyCost:  round2(totalWithBuffer),
	}
}

func round2(val float64) float64 {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}
	return math.Round(val*100) / 100
}
