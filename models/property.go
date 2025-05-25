package models

type Property struct {
	PropertyID int     `json:"property_id"`
	Type       string  `json:"type"`
	PAddress   string  `json:"p_address"`
	Prize      float64 `json:"prize"`
	MapLink    string  `json:"map_link"`
	Img        []byte  `json:"img_path"`
	CreatedAt  string  `json:"created_at"`
	UserID     int     `json:"user_id"`
}
