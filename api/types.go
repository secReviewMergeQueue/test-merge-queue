package train

import (
	pb "github.com/saumeya/train-ticketing/api/proto"
)

// Reservation represents a booking made by a user on a specific train route.
type Reservation struct {
	BookingID  string
	TrainID    string
	From       string
	To         string
	Price      float64
	User       pb.User
	SeatNumber string
}

// TrainRoute represents a specific train route including seating availability.
type TrainChart struct {
	TrainID  string
	From     string
	To       string
	Price    float64
	SectionA [20]bool // Availability of seats in Section A (20 seats).
	SectionB [20]bool // Availability of seats in Section B (20 seats).
}
