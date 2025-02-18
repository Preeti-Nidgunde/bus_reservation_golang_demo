package db

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB is the global database connection
var DB *sqlx.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return fmt.Errorf("Error: DB_PASSWORD environment variable is not set")
	}

	// Define DSN (Data Source Name) without hardcoding the password
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/", dbPassword)

	// Open the initial database connection
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("DB should not be initialized %v", err)
	}

	// Create the database if it doesn't exist
	_, err = DB.Exec("CREATE DATABASE IF NOT EXISTS bus_db;")
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// Update DSN to connect to the newly created `bus_db`
	dsn = fmt.Sprintf("root:%s@tcp(localhost:3306)/bus_db", dbPassword)

	// Connect to the `bus_db`
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to bus_db: %v", err)
	}

	// Verify the database connection
	// Verify the database connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("database is not reachable: %v", err)
	}

	// Call the function to create necessary tables
	// Call the function to create necessary tables
	createTables()

	return nil
}

func BeginTransaction(ctx context.Context) (*sqlx.Tx, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return DB.BeginTxx(ctx, nil)
}

func RollbackTranscation() {

}

func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS buses (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		seats INT NOT NULL
	);`,
		`CREATE TABLE IF NOT EXISTS seats (
		id INT AUTO_INCREMENT PRIMARY KEY,
		seat_number INT NOT NULL,
		status VARCHAR(20) DEFAULT 'available',
		user_id INT,
		bus_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
		FOREIGN KEY (bus_id) REFERENCES buses(id) ON DELETE CASCADE
	);`,
		`CREATE TABLE IF NOT EXISTS bookings (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		bus_id INT,
		seat_numbers JSON NOT NULL,
		booking_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (bus_id) REFERENCES buses(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range tables {
		_, err := DB.Exec(query)
		if err != nil {
			log.Fatal("Failed to create tables:", err)
		}
	}

	log.Println("Database tables created successfully")
}
