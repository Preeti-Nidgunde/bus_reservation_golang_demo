package bl

import (
	"bus_reservation/db"
	"bus_reservation/models"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tx, err := db.BeginTransaction(r.Context())
	if err != nil {
		http.Error(w, "Internal eroor", http.StatusInternalServerError)
		log.Println("Database error", err)
		return
	}

	// Ensure rollback if an error occurs
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		}
	}()

	if user.Username != "" && user.Password != "" {
		result, err := db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			log.Println("Failed to register a user : ", err)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println("Transaction commit failed:", err)
			http.Error(w, "Failed to register the user", http.StatusInternalServerError)
			return
		}

		// Get last inserted ID
		userID, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()

		log.Printf("User registered with ID: %d (Rows affected: %d)", userID, rowsAffected)
		log.Println("User registeration started")
	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Println("required fields are missing")
		return
	}

}

// Function to check if the user exists in the users table
func UserExists(userId int) bool {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id=?", userId).Scan(&count)
	if err != nil {
		log.Println("Error checking user existence:", err)
		return false
	}
	return count > 0
}
