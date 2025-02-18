package db

import (
	"context"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // blank import
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	DB, err = sqlx.Open("mysql", "root:qawzsx1@tcp(localhost:3306)/")
	if err != nil {
		log.Fatal("Failed to connect to MySQL DB")
	}
	defer DB.Close()

	log.Println("Connected to MYSQL Database !!")

	_, err = DB.Exec("CREATE DATABASE IF NOT EXISTS bus_db;")
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	DB, err = sqlx.Open("mysql", "root:qawzsx1@tcp(localhost:3306)/bus_db")
	if err != nil {
		log.Fatal("Failed to connect to railway_DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	createTables()

	log.Println("Created Bus database and required tables successfully !!")
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
