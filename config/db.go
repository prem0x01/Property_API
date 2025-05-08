package config

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
