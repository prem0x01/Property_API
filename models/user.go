package models

type User struct {
	UserID     int        `json:"user_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Mobile     string     `json:"mobile"`
	Password   string     `json:"password"`
	Aadhaar    int64      `json:"aadhaar"`
	UAddress   string     `json:"u_address"`
	UPFImg     []byte     `json:"upf_img_path"`
	Properties []Property `json:"properties"`
	CreatedAt  string     `json:"created_at"`
}
