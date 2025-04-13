package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// so we will perse all the html file throught the templates variable
var templates = template.Must(template.ParseGlob("templates/*.html"))

type Property struct {
	PropertyID int     `json:"property_id"`
	Type       string  `json:"type"`
	PAddress   string  `json:"p_address"`
	Prize      float64 `json:"prize"`
	MapLink    string  `json:"map_link"`
	ImgPath    string  `json:"img_path"`
}

type User struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Password   string `json:"password"`
	Aadhaar    int    `json:"aadhaar"`
	UAddress   string `json:"u_address"`
	UPFImgPath string `json:"upf_img_path"`
}

type Appointment struct {
	AppointmentID int    `json:"appointment_id"`
	UserID        int    `json:"user_id"`
	PropertyID    int    `json:"property_id"`
	Time          string `json:"time"`
	Date          string `json:"date"`
	Mobile        string `json:"mobile"`
	Address       string `json:"address"`
}

/*
so what is mutex?
mutex is used when we want to prevent go rutines to aceess same data simultenously
it prevent them by two method Lock() and Unlock() , lock loks the gorutine until it completes its task and
then unlock it after completion , it helps to privent race conditions
*/
var (
	mutex = &sync.Mutex{}
	db    *sql.DB
)

func main() {
	var err error

	connectDB()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripePrefix("/static/", fs))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/property", propertyHandler)
	http.HandleFunc("/appointment", appointmentHandler)

	fmt.Println("\033[35m[-] Starting server on :9090\033[0m")
	if err = http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("\033[31m[-] Could not start server: %s\n\033[0m", err.Error())
	}
}

func connectDB() {
	_, err := godotenv.Load()
	if err != nil {
		log.Fatalf("\033[31m[-] Can't open .env file: %v\n\033[0m ", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("\033[31m[-] Error opening database: %v\n\033[0m", err)
	}
	defer db.Close()

	//Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("\033[31m[-] Error connecting database: %v\n\033[0m", err)
	}

	fmt.Println("\033[35m[-] Connected to database successfully!\033[0m")

	createTables()
}
func createTables() {

	createUserTable := `
	CREATE TABLE users (
    	user_id SERIAL PRIMARY KEY,  -- Auto-incrementing primary key
		name TEXT NOT NULL
    	email VARCHAR(255) UNIQUE NOT NULL,  -- Unique email (required)
    	mobile VARCHAR(15) UNIQUE NOT NULL,  -- Mobile number (required)
    	password TEXT NOT NULL,  -- Hashed password storage
    	aadhaar BIGINT UNIQUE NOT NULL,  -- Aadhaar number (must be unique)
		u_address TEXT,  -- User address (optional)
    	upf_img_path TEXT,  -- Profile image path (optional)
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Auto-set timestamp on creation
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- Auto-update timestamp on modification
	);`

	createPropertyTable := `
	CREATE TABLE properties (
    	property_id SERIAL PRIMARY KEY,  -- Auto-incrementing unique ID
    	type VARCHAR(50) NOT NULL,       -- Property type (e.g., Apartment, Villa)
    	p_address TEXT NOT NULL,         -- Property address (required)
    	prize DECIMAL(12,2) NOT NULL,    -- Price with precision for monetary values
    	map_link TEXT,                   -- Google Maps link (optional)
    	img_path TEXT,                    -- Image file path (optional)
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Auto-set timestamp on creation
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Auto-update timestamp on modification
	);`

	createAppointmentTable := `
	CREATE TABLE appointments (
    	appointment_id SERIAL PRIMARY KEY, -- Auto-incrementing unique ID
    	user_id INT REFERENCES users(user_id) ON DELETE CASCADE, -- Linked to users table
    	property_id INT REFERENCES properties(property_id) ON DELETE CASCADE, -- Linked to properties table
    	time TIME NOT NULL, -- Stores appointment time
    	date DATE NOT NULL, -- Stores appointment date
    	mobile VARCHAR(15) NOT NULL, -- Contact number (required)
    	address TEXT NOT NULL, -- Location or meeting place
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Track creation time
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Track modifications
	);`

	if _, err := db.Exec(createUserTable); err != nil {
		log.Fatalf("\033[31m[-] Error creating users table: %v\n\033[0m", err)
	}
	if _, err := db.Exec(createPropertyTable); err != nil {
		log.Fatalf("\033[31m[-] Error creating propertys table: %v\n\033[0m", err)
	}
	if _, err := db.Exec(createAppointmentTable); err != nil {
		log.Fatalf("\033[31m[-] Error creating appointments table: %v\n\033[0m", err)
	}

	fmt.Println("\033[35m[-] Table Created Successfully\033[0m")

}

func userHandler(w http.ResponseWriter, r *http.Request) {
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

// now I have to create this crud functions , 12 APR 2025 4:47
// Rest i will do tomorowwwwwwwwww...

func isValidAadhaar(aadhaar string) bool {
	re := regexp.MustCompile(`^[0-9]{12}$`)
	return re.MatchString(aadhaar)
}

func isValidMobile(mobile string) bool {
	re := regexp.MustCompile(`[0-9]{10}$`)
	return re.MatchString(mobile)
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

	var result []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email, &u.Mobile, &u.Aadhaar, &u.UAddress); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func addUser(w http.ResponseWriter,  r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidAadhaar(u.Aadhaar) || !isValidMobile(u.Mobile) {
		http.Error(w, "Invalid Addhar or Mobile number format", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Query("INSERT INTO users(name, email, mobile , password, aadhaar, u_address, upf_img_path ) VALUES(?, ?, ?, ?, ?, ?, ?)")
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
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http..Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidAadhaar(u.UAddress) || !isValidMobile(u.Mobile) {
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
	json.NewDecoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request){
	idStr := r.URL.Query().Get("user_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	stmt, err := db.Prepare("DELETE * FROM users WHERE user_id = ?" )
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Today i wrote all users CRUD functions . 13 APR 2025 5.00
