package main

import (
	"log"
	"net/http"

	"bus_reservation/db"
	"bus_reservation/endpoints"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	router := mux.NewRouter()
	endpoints.AddHandlers(router)
	log.Fatal(http.ListenAndServe(":8090", router))

}
