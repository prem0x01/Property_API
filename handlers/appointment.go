package handlers

import (
	"Property_App/models"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var (
	db    *sql.DB
	mutex = &sync.Mutex{}
)

func InitAppointmentHandler(database *sql.DB) {
	db = database
}

func AppointmentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewAppointment(w, r)
	case "POST":
		addAppointment(w, r)
	case "PUT":
		updateAppointment(w, r)
	case "DELETE":
		deleteAppointment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func viewAppointment(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	rows, err := db.Query(`SELECT 
		a.appointment_id, a.time, a.date, a.mobile, a.address, u.user_id, u.name, u.email,
		p.property_id, p.type, p.p_address, p.prize, p.map_link, p.image
		FROM appointments a
		JOIN users u ON a.user_id = u.user_id
		JOIN properties p ON a.property_id = p.property_id`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var appointments []struct {
		Appointment models.Appointment `json:"appointment"`
		UserName    string             `json:"user_name"`
		UserEmail   string             `json:"user_email"`
		Property    models.Property    `json:"property"`
	}

	for rows.Next() {
		var a models.Appointment
		var p models.Property
		var imageData []byte
		var userName, userEmail string

		if err := rows.Scan(&a.AppointmentID, &a.Time, &a.Date, &a.Mobile, &a.Address,
			&a.UserID, &userName, &userEmail,
			&p.PropertyID, &p.Type, &p.PAddress, &p.Prize, &p.MapLink, &imageData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p.Img = base64.StdEncoding.EncodeToString(imageData)

		appointments = append(appointments, struct {
			Appointment models.Appointment `json:"appointment"`
			UserName    string             `json:"user_name"`
			UserEmail   string             `json:"user_email"`
			Property    models.Property    `json:"property"`
		}{
			Appointment: a,
			UserName:    userName,
			UserEmail:   userEmail,
			Property:    p,
		})

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

func addAppointment(w http.ResponseWriter, r *http.Request) {
	var a models.Appointment
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("INSERT INTO appointments(user_id, property_id, time, date, mobile, address) VALUES($1, $2, $3, $4, $5, $6) RETURNING appointment_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(a.UserID, a.PropertyID, a.Time, a.Date, a.Mobile, a.Address).Scan(&a.AppointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func updateAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	appointmentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}
	var a models.Appointment
	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("UPDATE appointments SET time=$1, date=$2, mobile=$3, address=$4 WHERE appointment_id=$5")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Time, a.Date, a.Mobile, a.Address, appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No appointment found with the given ID", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Appointment updated successfully"})
}

func deleteAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	appointmentID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE FROM appointments WHERE appointment_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(appointmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No appointment found with the given ID", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
