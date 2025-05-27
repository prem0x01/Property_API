package models

type User struct {
	UserID       int        `json:"user_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Mobile       string     `json:"mobile"`
	Password     string     `json:"password"`
	Aadhaar      int64      `json:"aadhaar"`
	UAddress     string     `json:"u_address"`
	UPFImg       []byte     `json:"-"`
	UPFImgBase64 string     `json:"upf_img,omitempty"`
	Properties   []Property `json:"properties,omitempty"`
	CreatedAt    string     `json:"created_at"`
}
