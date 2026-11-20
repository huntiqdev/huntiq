package dto

// LocationQuery holds query params for the location hierarchy endpoints.
// Used by ListZones, ListSubcities, ListWoredas, ListKebeles, and SearchLocations.
type LocationQuery struct {
	Region   string   `query:"region"`
	Zone     string   `query:"zone"`
	Subcity  string   `query:"subcity"`
	Woreda   string   `query:"woreda"`
	Kebele   string   `query:"kebele"`
	Q        string   `query:"q"`         // search term for SearchLocations
	Lat      *float64 `query:"lat"`       // GPS latitude coordinate (free for all users)
	Lng      *float64 `query:"lng"`       // GPS longitude coordinate (free for all users)
	RadiusKM *float64 `query:"radius_km"` // search radius in km for nearby listings
}
