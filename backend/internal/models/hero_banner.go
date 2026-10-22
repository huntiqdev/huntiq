package models

import "time"

type HeroBanner struct {
	ID             int        `json:"id"`
	BannerType     string     `json:"banner_type"`     // "buyer" | "seller"
	MediaType      string     `json:"media_type"`      // "image" | "video"
	MediaURL       string     `json:"media_url"`
	Title          string     `json:"title"`
	Subtitle       string     `json:"subtitle"`
	LinkURL        string     `json:"link_url"`
	IsActive       bool       `json:"is_active"`
	SortOrder      int        `json:"sort_order"`
	PromotionText  string     `json:"promotion_text,omitempty"`  // overlay promo message
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
	TargetAudience string     `json:"target_audience"` // all|buyer|seller|normal|gold|premium
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
