package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("\033[31m[-] Can't open .env file: %v\n\033[0m ", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("\033[31m[-] Error opening database: %v\n\033[0m", err)
	}
	//defer db.Close()  causing race condition , its closing the connection befor running CreateTables(), place db.Close in main.

	err = db.Ping()
	if err != nil {
		log.Fatalf("\033[31m[-] Error connecting database: %v\n\033[0m", err)
	}

	fmt.Println("\033[35m[-] Connected to database successfully!\033[0m")

	createTables(db)
	return db, nil
}
func createTables(db *sql.DB) {

	createUserTable := `
	CREATE TABLE IF NOT EXISTS user (
    	user_id SERIAL PRIMARY KEY,
    	name TEXT NOT NULL,
    	email VARCHAR(255) UNIQUE NOT NULL,
		mobile VARCHAR(15) UNIQUE NOT NULL,
    	password TEXT NOT NULL,
    	aadhaar BIGINT UNIQUE NOT NULL,
    	u_address TEXT,
    	upf_img BYTEA,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createPropertyTable := `
	CREATE TABLE IF NOT EXISTS properties (
    	property_id SERIAL PRIMARY KEY,  
    	type VARCHAR(50) NOT NULL,      
    	p_address TEXT NOT NULL,         
    	prize DECIMAL(12,2) NOT NULL,    
    	map_link TEXT,                   
    	img BYTEA,                  
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createAppointmentTable := `
	CREATE TABLE IF NOT EXISTS appointments (
    	appointment_id SERIAL PRIMARY KEY, 
    	user_id INT REFERENCES users(user_id) ON DELETE CASCADE, 
    	property_id INT REFERENCES properties(property_id) ON DELETE CASCADE, 
    	time TIME NOT NULL, 
    	date DATE NOT NULL, 
    	mobile VARCHAR(15) NOT NULL, 
    	address TEXT NOT NULL, 
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
