package handlers

import (
	//"Property_App/config"
	"Property_App/models"
	//"Property_App/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
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

	rows, err := db.Query("SELECT user_id, appointment_id, property_id, time, date, mobile, address FROM appointments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.UserID, &a.AppointmentID, &a.PropertyID, &a.Time, &a.Date, &a.Mobile, &a.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
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

	stmt, err := db.Prepare("INSERT INTO appointments(user_id, appointment_id, property_id, time, date, mobile, address) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(&a.UserID, &a.AppointmentID, &a.PropertyID, &a.Time, &a.Date, &a.Mobile, &a.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.PropertyID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func updateAppointment(w http.ResponseWriter, r *http.Request) {
	var a models.Appointment
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("UPDATE appointments SET time=?, date=?, mobile=?, address=? WHERE appointment_id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Time, a.Date, a.Mobile, a.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

func deleteAppointment(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("appointment_id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE * FROM appointments WHERE appointment_id = ?")
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
