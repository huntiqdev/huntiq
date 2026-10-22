package dto

// LocationQuery holds query params for the location hierarchy endpoints.
// Used by ListZones, ListSubcities, ListWoredas, and SearchLocations.
type LocationQuery struct {
	Region  string `query:"region"`
	Zone    string `query:"zone"`
	Subcity string `query:"subcity"`
	Q       string `query:"q"` // search term for SearchLocations
}
