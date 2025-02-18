package bl

import (
	"bus_reservation/db"
	"bus_reservation/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var mu sync.Mutex

func FreeSeats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	busIDStr := vars["bus_id"]
	busID, err := strconv.Atoi(busIDStr)
	if err != nil {
		http.Error(w, "Invalid bus ID", http.StatusBadRequest)
		log.Println("Invalid bus ID")
		return
	}

	ctx := r.Context()

	rows, err := db.DB.QueryContext(ctx, "SELECT seat_number FROM seats WHERE bus_id = ? AND status = 'available'", busID)
	if err != nil {
		http.Error(w, "Failed to get free seats", http.StatusInternalServerError)
		log.Println("Failed to get free seats")
		return
	}
	defer rows.Close()

	var freeSeats []int
	for rows.Next() {
		var seatNumber int
		if err := rows.Scan(&seatNumber); err != nil {
			http.Error(w, "Failed to scan seat number", http.StatusInternalServerError)
			return
		}
		freeSeats = append(freeSeats, seatNumber)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over seats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(freeSeats)
}

func createBus(ctx context.Context, bus models.Bus) (int, error) {
	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %v", err)
	}

	// Ensure rollback if an error occurs
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		}
	}()

	if bus.Seats <= 0 {
		log.Println("Invalid Bus seats")
		return 0, fmt.Errorf("Invalid Bus seats")
	}

	// Insert bus using transaction
	result, err := tx.Exec("INSERT INTO buses (Name, Seats) VALUES (?, ?)", bus.Name, bus.Seats)
	if err != nil {
		log.Println("Error executing query:", err)
		return 0, fmt.Errorf("Failed to create a given bus: %v", err)
	}

	// Get last inserted bus ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting LastInsertId:", err)
		return 0, fmt.Errorf("Failed to retrieve the last inserted ID: %v", err)
	}

	log.Println("Bus created with ID:", lastInsertID)

	// Insert seats using transaction
	for i := 1; i <= bus.Seats; i++ {
		_, err := tx.Exec("INSERT INTO seats (bus_id, seat_number, status) VALUES (?, ?, 'available')", lastInsertID, i)
		if err != nil {
			log.Printf("Failed to insert seat %d for bus %d: %v", i, lastInsertID, err)
			return 0, fmt.Errorf("failed to insert seat %d: %v", i, err)
		}
	}

	// Commit transaction if all operations succeed
	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit failed:", err)
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Transaction committed successfully")
	return int(lastInsertID), nil
}

func CreateBusHandler(w http.ResponseWriter, r *http.Request) {
	var bus models.Bus
	var err error

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&bus); err != nil {
		http.Error(w, "Unable to parse the request body", http.StatusBadRequest)
		log.Println("Error decoding JSON:", err)
		return
	}

	ctx := r.Context()
	bus.Id, err = createBus(ctx, bus)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating bus: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"Message": "Bus Added successfully", "BusID": bus.Id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetAllBuses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("Select * from buses")
	if err != nil {
		http.Error(w, "Unable to fetch the requested data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var buses []models.Bus
	for rows.Next() {
		var bus models.Bus
		if err := rows.Scan(&bus.Id, &bus.Name, &bus.Seats); err != nil {
			http.Error(w, "Failed to parse the buses", http.StatusInternalServerError)
		}
		buses = append(buses, bus)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buses)
}

func BusExists(busId int) bool {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM BUSES WHERE id=?", busId).Scan(&count)
	if err != nil {
		log.Println("Error checking bus existence:", err)
		return false
	}
	return count > 0
}

func ValidateBusId(r *http.Request) int {
	vars := mux.Vars(r)
	busIdstr, ok := vars["bus_id"]
	if !ok {
		log.Println("Bus ID missing in URL")
		return -1
	}

	busId, err := strconv.Atoi(busIdstr)
	if err != nil {
		log.Println("Invalid value for BusID", err)
		return -1
	}

	if !BusExists(busId) {
		return -1
	}
	return busId
}
