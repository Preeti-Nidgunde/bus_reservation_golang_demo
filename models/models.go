package models

// User struct
type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Train struct
type Bus struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Seats int    `json:"seats"`
}

// Booking Struct
type Booking struct {
	Id          int   `json:"id"`
	UserId      int   `json:"user_id"`
	BusId       int   `json:"bus_id"`
	SeatNumbers []int `json:"seat_numbers"`
}

type Seats struct {
	Id     int    `json:"id"`
	BusId  int    `json:"bus_id"`
	UserId int    `json:"user_id"`
	Status string `json: status`
}
