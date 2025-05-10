package handlers

import (
	"Property_App/models"
	"strconv"
	//"Property_App/utils"
	//"Property_App/config"
	"database/sql"
	"encoding/json"
	"net/http"
	//"sync"
)

func InitPropertyHandler(database *sql.DB) {
	db = database
}

func PropertyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewProperties(w, r)
	case "POST":
		addProperty(w, r)
	case "PUT":
		updateProperty(w, r)
	case "DELETE":
		deleteProperty(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func viewProperties(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query("SELECT property_id, type, p_address, prize, map_link, img_path FROM properties")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var properties []models.Property
	for rows.Next() {
		var p models.Property
		if err := rows.Scan(&p.PropertyID, &p.Type, &p.PAddress, &p.Prize, &p.MapLink, &p.ImgPath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		properties = append(properties, p)
	}
	json.NewEncoder(w).Encode(properties)

}

func addProperty(w http.ResponseWriter, r *http.Request) {
	var p models.Property
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("INSERT INTO properties(type ,p_address, prize, map_link, img_path) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Type, p.PAddress, p.Prize, p.MapLink, p.ImgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.PropertyID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updateProperty(w http.ResponseWriter, r *http.Request) {
	var p models.Property
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("UPDATE properties SET type=?, p_address=?, prize=?, map_link=?, img_path=? WHERE property_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Type, p.PAddress, p.Prize, p.MapLink, p.ImgPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deleteProperty(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("property_id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE * FROM properties WHERE property_id = ?")
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
