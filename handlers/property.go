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

func InitPropertyHandler(database *sql.DB) {
	db = database
}

func propertyHandler(w http.ResponseWriter, r *http.Request) {
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
