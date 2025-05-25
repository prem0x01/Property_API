package main

import (
	"Property_App/config"
	"Property_App/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("\nFailed to connect to the database: %v", err)
	}
	defer db.Close()

	//config.CreateTables()

	handlers.InitUserHandler(db)
	handlers.InitPropertyHandler(db)
	handlers.InitAppointmentHandler(db)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/property", handlers.PropertyHandler)
	http.HandleFunc("/appointment", handlers.AppointmentHandler)

	fmt.Println("\033[35m[-] Server running on :9090....\033[0m")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
