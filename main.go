package main

import (
	"Property_App/config"
	"Property_App/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("\033[31m[-] Could not load .env file: %v\n\033[0m", err)
	}

	db := config.InitDB()
	defer db.Close()

	config.CreateTables(db)

	handlers.InitUserHandler(db)
	handlers.InitPropertyHandler(db)
	handlers.InitAppointmentHandler(db)

	// serveing static files . thats new thing i leanr about StripPrefix()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/property", handlers.PropertyHandler)
	http.HandleFunc("/appointment", handlers.AppointmentHandler)

	fmt.Println("\033[35m[-] Server running on :9090....\033[0m")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
