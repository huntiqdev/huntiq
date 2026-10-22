package models

import "time"

type Location struct {
	ID        int       `json:"id"`
	Region    string    `json:"region"`
	Zone      *string   `json:"zone"`
	Subcity   *string   `json:"subcity"`
	Woreda    *string   `json:"woreda"`
	RegionAm  *string   `json:"region_am,omitempty"`
	RegionOr  *string   `json:"region_or,omitempty"`
	RegionTi  *string   `json:"region_ti,omitempty"`
	ZoneAm    *string   `json:"zone_am,omitempty"`
	ZoneOr    *string   `json:"zone_or,omitempty"`
	ZoneTi    *string   `json:"zone_ti,omitempty"`
	SubcityAm *string   `json:"subcity_am,omitempty"`
	SubcityOr *string   `json:"subcity_or,omitempty"`
	SubcityTi *string   `json:"subcity_ti,omitempty"`
	WoredaAm  *string   `json:"woreda_am,omitempty"`
	WoredaOr  *string   `json:"woreda_or,omitempty"`
	WoredaTi  *string   `json:"woreda_ti,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

