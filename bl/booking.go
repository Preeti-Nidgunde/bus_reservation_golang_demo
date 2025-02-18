package bl

import (
	"bus_reservation/db"
	"bus_reservation/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Get all available seats for a bus
func GetAvailableSeats(w http.ResponseWriter, r *http.Request) {
	busId := ValidateBusId(r)
	if busId == -1 {
		http.Error(w, "Incorrect Bus ID", http.StatusNotFound)
		log.Println("Incorrect Bus ID")
		return
	}

	freeSeats := []int{}
	rows, err := db.DB.Query("Select seat_number from seats where status = 'available' AND bus_id=?", busId)
	if err != nil {
		http.Error(w, "Unable to fetch the seats", http.StatusInternalServerError)
		log.Println("Unable to fetch the seats", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var seatNo int
		if err := rows.Scan(&seatNo); err != nil {
			http.Error(w, "Failed to parse the buses", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		freeSeats = append(freeSeats, seatNo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(freeSeats)
}

func TotalBookings(w http.ResponseWriter, r *http.Request) {
	var totalBookings int
	busId := ValidateBusId(r)
	if busId == -1 {
		http.Error(w, "Incorrect Bus ID", http.StatusBadRequest)
		log.Println("Incorrect Bus ID")
		return
	}

	err := db.DB.QueryRow("SELECT COUNT(*) FROM SEATS WHERE STATUS!='available' AND bus_id=?", busId).Scan(&totalBookings)
	if err != nil {
		http.Error(w, "Unable to get the count", http.StatusInternalServerError)
		log.Println("Unable to to get the count", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"Count of booked seats": totalBookings})
}

func bookSeats(ctx context.Context, booking models.Booking) (int, []int, error) {
	if !UserExists(booking.UserId) {
		return 0, nil, fmt.Errorf("user does not exist, please register first")
	}

	tx, err := db.BeginTransaction(ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Println("Transaction rolled back due to error:", err)
		}
	}()

	var bookedSeats []int
	var unbookedSeats []int

	// Check availability of each seat before booking
	for _, seatNumber := range booking.SeatNumbers {
		log.Print("seatNumber = ", seatNumber)
		var seatStatus string
		err := tx.QueryRow("SELECT status FROM seats WHERE seat_number = ? AND bus_id = ?", seatNumber, booking.BusId).Scan(&seatStatus)
		if err != nil {
			log.Printf("Error checking seat %d: %v", seatNumber, err)
			return 0, nil, fmt.Errorf("failed to check availability for seat %d: %v", seatNumber, err)
		}

		if seatStatus != "available" {
			log.Printf("Seat %d is already booked for bus %d", seatNumber, booking.BusId)
			unbookedSeats = append(unbookedSeats, seatNumber)
			continue
		}

		_, err = tx.Exec("UPDATE seats SET status = 'booked', user_id = ? WHERE seat_number = ? AND bus_id = ?", booking.UserId, seatNumber, booking.BusId)
		if err != nil {
			log.Printf("Failed to update seat %d for bus %d: %v", seatNumber, booking.BusId, err)
			return 0, nil, fmt.Errorf("failed to book seat %d: %v", seatNumber, err)
		}
		bookedSeats = append(bookedSeats, seatNumber)
	}

	bookedSeatsJson, err := json.Marshal(bookedSeats) // Convert the booked seats array to JSON format
	if err != nil {
		log.Printf("Error marshalling booked seats: %v", err)
		return 0, nil, fmt.Errorf("failed to marshal booked seats: %v", err)
	}
	_, err = tx.Exec("INSERT INTO BOOKINGS (user_id, bus_id, seat_numbers) VALUES (?, ?, ?)", booking.UserId, booking.BusId, bookedSeatsJson)
	if err != nil {
		log.Printf("Failed to insert booking record for bus %d: %v", booking.BusId, err)
		return 0, nil, fmt.Errorf("failed to create booking: %v", err)
	}

	// Commit the transaction if all operations succeed
	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit failed:", err)
		return 0, nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("Transaction committed successfully")
	return len(bookedSeats), unbookedSeats, nil
}

func BookSeatsHandler(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Unable to parse the request", http.StatusBadRequest)
		log.Println("Bad request", err)
		return
	}

	busId := booking.BusId
	if !BusExists(busId) {
		http.Error(w, "Incorrect Bus ID", http.StatusNotFound)
		log.Println("Incorrect Bus ID")
		return
	}

	ctx := r.Context()
	bookedCount, unbookedSeats, err := bookSeats(ctx, booking)
	if err != nil {
		// Handle user not found error
		if err.Error() == "user does not exist, please register first" {
			http.Error(w, "User does not exist, please register first", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Error creating booking: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare response message
	response := map[string]interface{}{
		"Message":       "Booking done successfully",
		"BookedSeats":   bookedCount,
		"UnbookedSeats": len(unbookedSeats),
	}

	// Set response header and send back the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
