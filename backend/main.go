package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	m "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Load .env file
	err := godotenv.Load(".env." + env + ".local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()

	// Database connection setup
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("Can't open database AB")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("Pinging failed")
		log.Fatal(pingErr)
	}
	log.Println("Connected to the database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Migrations")

	r := mux.NewRouter()

	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/users", getUsersHandler).Methods("GET")

	frontend := os.Getenv("HOST")
	frontend_port := os.Getenv("HOST_PORT")

	log.Printf("frontend: http://%s:%s", frontend, frontend_port)

	// handle Cors
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://%s:%s", frontend, frontend_port},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{"Content-Type", "application/json"},
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func runMigrations(db *sql.DB) error {
	driver, err := m.WithInstance(db, &m.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	result, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting user ID", http.StatusInternalServerError)
		return
	}

	user.ID = int(id)
	user.Password = "" // Don't send the password back

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginAttempt User
	err := json.NewDecoder(r.Body).Decode(&loginAttempt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", loginAttempt.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	print(&user.Email)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginAttempt.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	user.Password = "" // Don't send the password back

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user":    user,
	})
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, email FROM users")
	if err != nil {
		log.Println(err)

	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Println(users)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
