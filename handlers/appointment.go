package handlers

import (
	"Property_App/models"
	"Property_App/utils"
	"database/sql"
	"sync"
)

var (
	db    *sql.DB
	mutex = &sync.Mutex{}
)

func InitAppointmentHandler(database *sql.DB) {
	db = database
}

func appointmentHandler(w http.ResponseWriter, r *http.Request) {
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
