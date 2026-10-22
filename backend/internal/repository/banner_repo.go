package repository

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/used-tools-marketplace/internal/models"
	"github.com/your-org/used-tools-marketplace/pkg/constants"
)

type BannerRepo struct{ db *pgxpool.Pool }

func NewBannerRepo(db *pgxpool.Pool) *BannerRepo { return &BannerRepo{db: db} }

// maxBannersForType reads the admin-configured max active banners for a banner type from platform_settings.
func (r *BannerRepo) maxBannersForType(ctx context.Context, bannerType string) int {
	key := constants.SettingMaxBannersBuyer
	if bannerType == "seller" {
		key = constants.SettingMaxBannersSeller
	}
	var val string
	_ = r.db.QueryRow(ctx, `SELECT value FROM platform_settings WHERE key=$1`, key).Scan(&val) //nolint:gosec
	if val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			return n
		}
	}
	return constants.DefaultMaxBanners
}

const bannerListActiveSQL = `SELECT id, banner_type, COALESCE(media_type,'image'), media_url,
		        COALESCE(title,''), COALESCE(subtitle,''), COALESCE(link_url,''),
		        is_active, sort_order, created_at, updated_at
		 FROM hero_banners
		 WHERE banner_type=$1 AND is_active=true ORDER BY sort_order, id LIMIT $2`

const bannerListAllSQL = `SELECT id, banner_type, COALESCE(media_type,'image'), media_url,
		        COALESCE(title,''), COALESCE(subtitle,''), COALESCE(link_url,''),
		        is_active, sort_order, created_at, updated_at
		 FROM hero_banners ORDER BY banner_type, sort_order, id`

func (r *BannerRepo) ListActive(ctx context.Context, bannerType string) ([]models.HeroBanner, error) {
	limit := r.maxBannersForType(ctx, bannerType)
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, bannerListActiveSQL, bannerType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBanners(rows)
}

func (r *BannerRepo) ListAll(ctx context.Context) ([]models.HeroBanner, error) {
	//nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	rows, err := r.db.Query(ctx, bannerListAllSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBanners(rows)
}

func (r *BannerRepo) Create(ctx context.Context, b models.HeroBanner) (int, error) {
	if b.MediaType == "" {
		b.MediaType = "image"
	}
	var id int
	err := r.db.QueryRow(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`INSERT INTO hero_banners (banner_type, media_type, media_url, title, subtitle, link_url, is_active, sort_order)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		b.BannerType, b.MediaType, b.MediaURL, b.Title, b.Subtitle, b.LinkURL, b.IsActive, b.SortOrder,
	).Scan(&id)
	return id, err
}

func (r *BannerRepo) Update(ctx context.Context, b models.HeroBanner) error {
	if b.MediaType == "" {
		b.MediaType = "image"
	}
	_, err := r.db.Exec(ctx, //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
		`UPDATE hero_banners SET banner_type=$1, media_type=$2, media_url=$3, title=$4,
		 subtitle=$5, link_url=$6, is_active=$7, sort_order=$8, updated_at=NOW()
		 WHERE id=$9`,
		b.BannerType, b.MediaType, b.MediaURL, b.Title, b.Subtitle, b.LinkURL, b.IsActive, b.SortOrder, b.ID)
	return err
}

func (r *BannerRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM hero_banners WHERE id=$1`, id) //nolint:gosec // G201/G202: all values passed as $N parameters, no user input in SQL structure
	return err
}

type bannerRows interface {
	Next() bool
	Scan(...any) error
	Err() error
}

func scanBanners(rows bannerRows) ([]models.HeroBanner, error) {
	var out []models.HeroBanner
	for rows.Next() {
		var b models.HeroBanner
		if err := rows.Scan(&b.ID, &b.BannerType, &b.MediaType, &b.MediaURL,
			&b.Title, &b.Subtitle, &b.LinkURL, &b.IsActive, &b.SortOrder,
			&b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
