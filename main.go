package main

import (
	"Property_App/config"
	"Property_App/handlers"
	"Property_App/utils"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
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
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.Handle("/user", utils.RateLimiter(http.HandlerFunc(handlers.UserHandler))).Methods("GET", "POST")
	router.Handle("/user/{id}", utils.RateLimiter(http.HandlerFunc(handlers.UserHandler))).Methods("DELETE", "PUT")
   
   router.Handle("/property", utils.RateLimiter(http.HandlerFunc(handlers.PropertyHandler))).Methods("GET", "POST")
	router.Handle("/property/{id}", utils.RateLimiter(http.HandlerFunc(handlers.PropertyHandler))).Methods("DELETE", "PUT")

	router.Handle("/appointment", utils.RateLimiter(http.HandlerFunc(handlers.AppointmentHandler))).Methods("GET", "POST")
	router.Handle("/appointment/{id}", utils.RateLimiter(http.HandlerFunc(handlers.AppointmentHandler))).Methods("DELETE", "PUT")

	fmt.Println("\033[35m[-] Server running on :9090....\033[0m")
	log.Fatal(http.ListenAndServe("localhost:9090", router))

}
