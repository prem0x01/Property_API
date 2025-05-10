package handlers

import (
	"Property_App/models"
	"Property_App/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	//"sync"
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
	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query("SELECT user_id, name, email, mobile, aadhaar, u_address FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Mobile, &u.Aadhaar, &u.UAddress); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.IsValidAadhaar(u.Aadhaar) || !utils.IsValidMobile(u.Mobile) {
		http.Error(w, "Invalid Addhar or Mobile number format", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("INSERT INTO users(name, email, mobile , password, aadhaar, u_address, upf_img_path ) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Name, u.Email, u.Mobile, u.Password, u.Aadhaar, u.UAddress, u.UPFImgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.UserID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.IsValidAadhaar(u.Aadhaar) || !utils.IsValidMobile(u.Mobile) {
		http.Error(w, "Invalid Aadhaar or Mobile number format", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("UPDATE users SET name=?, email=?, mobile=?, password=?, aadhaar=?, u_address=?, upf_img_path=? WHERE user_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Email, u.Mobile, u.Password, u.Aadhaar, u.UAddress, u.UPFImgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("user_id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE * FROM users WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
