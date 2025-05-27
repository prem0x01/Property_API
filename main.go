package main

import (
	"Property_App/config"
	"Property_App/handlers"
	"Property_App/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.InitLogger()
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

	http.Handle("/user", utils.RateLimiter(http.HandlerFunc(handlers.UserHandler)))
	http.Handle("/property", utils.RateLimiter(http.HandlerFunc(handlers.PropertyHandler)))
	http.Handle("/appointment", utils.RateLimiter(http.HandlerFunc(handlers.AppointmentHandler)))

	fmt.Println("\033[35m[-] Server running on :9090....\033[0m")
	log.Fatal(http.ListenAndServe("localhost:9090", nil))

}
