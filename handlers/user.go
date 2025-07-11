package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prem0x01/propertyAPI/models"
	"github.com/prem0x01/propertyAPI/utils"
)

func InitUserHandler(database *sql.DB) {
	db = database
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewUser(w, r)
	case "POST":
		addUser(w, r)
	case "PUT":
		updateUser(w, r)
	case "DELETE":
		deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func viewUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query(`
		SELECT
			u.user_id, u.name, u.email, u.mobile, u.aadhaar, u.u_address, u.upf_img,
			p.property_id, p.type, p.p_address, p.prize, p.map_link, p.img
		FROM users u
		LEFT JOIN properties p ON u.user_id = p.user_id
		WHERE u.user_id = $1
	`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var user models.User
	var properties []models.Property
	var imageData []byte

	for rows.Next() {
		var property models.Property
		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Mobile, &user.Aadhaar, &user.UAddress, &imageData,
			&property.PropertyID, &property.Type, &property.PAddress, &property.Prize, &property.MapLink, &property.Img); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.UPFImg = imageData
		user.UPFImgBase64 = base64.StdEncoding.EncodeToString(imageData)

		if property.PropertyID != 0 {
			properties = append(properties, property)
		}
	}

	user.Properties = properties

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(15 << 20) // Max 20MB file size
	if err != nil {
		http.Error(w, "File too large or invalid form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("upf_img")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileBytes []byte
	if file != nil {
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
	}

	var u models.User

	u.Name = r.FormValue("name")
	u.Email = r.FormValue("email")
	u.Mobile = r.FormValue("mobile")
	u.Password = r.FormValue("password")
	u.Aadhaar, _ = strconv.ParseInt(r.FormValue("aadhaar"), 10, 64)
	u.UAddress = r.FormValue("u_address")
	u.UPFImg = fileBytes

	hashPass, err := utils.HashPassword(u.Password)
	if err != nil {
		return
	}

	// Validate Aadhaar and Mobile
	if !utils.IsValidAadhaar(u.Aadhaar) || !utils.IsValidMobile(u.Mobile) {
		http.Error(w, "Invalid Aadhaar or Mobile number format", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Use QueryRow().Scan() with RETURNING to get inserted user_id
	err = db.QueryRow(`
		INSERT INTO users (name, email, mobile, password, aadhaar, u_address, upf_img)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id
	`, u.Name, u.Email, u.Mobile, hashPass, u.Aadhaar, u.UAddress, u.UPFImg).Scan(&u.UserID)

	if err != nil {
		http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.UserID == 0 {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if !utils.IsValidAadhaar(u.Aadhaar) || !utils.IsValidMobile(u.Mobile) {
		http.Error(w, "Invalid Aadhaar or Mobile number format", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare(`UPDATE users
		SET name=$1, email=$2, mobile=$3, password=$4, aadhaar=$5, u_address=$6, upf_img=$7
		WHERE user_id=$8`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Name, u.Email, u.Mobile, u.Password, u.Aadhaar, u.UAddress, u.UPFImg, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE FROM users WHERE user_id = $1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
