package tests

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
)

func getUserProfileImage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var imageData []byte
		err := db.QueryRow("SELECT upf_img FROM users WHERE user_id = $1", r.URL.Query().Get("user_id")).Scan(&imageData)
		if err != nil {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}

		base64Img := base64.StdEncoding.EncodeToString(imageData)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"image": "data:image/jpeg;base64,%s"}`, base64Img)
	}
}
