package endpoints

import (
	"bus_reservation/bl"
	"github.com/gorilla/mux"
	"log"
)

func AddHandlers(router *mux.Router) {
	router.HandleFunc("/register", bl.RegisterUser).Methods("POST")
	router.HandleFunc("/buses", bl.GetAllBuses).Methods("GET")
	router.HandleFunc("/bus", bl.CreateBusHandler).Methods("POST")
	router.HandleFunc("/availableSeats/{bus_id}", bl.GetAvailableSeats).Methods("GET")
	router.HandleFunc("/totalBookings/{bus_id}", bl.TotalBookings).Methods("GET")
	router.HandleFunc("/bookseats", bl.BookSeatsHandler).Methods("POST")
	log.Println("Routes registered using Mux router")
}
