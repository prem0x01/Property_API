package models

type User struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Password   string `json:"password"`
	Aadhaar    int    `json:"aadhaar"`
	UAddress   string `json:"u_address"`
	UPFImgPath string `json:"upf_img_path"`
}
