package bookingservice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	train "github.com/saumeya/train-ticketing/api"
)

// this function returns the train details for the requested route From source To destination along with unit price of the ticket
func (s *BookingServer) getTrainDetails(from string, to string) (train_id string, price float64, err error) {

	// checking if the train route exists
	for _, train := range *s.trainChart {
		if train.From == from && train.To == to {
			return train.TrainID, train.Price, nil
		}
	}

	// No train found, default return values are used
	err = errors.New("no trains with available seats for given route")
	return

}

// this function allocates the seat to the user requesting a ticket for a particular TrainID
// the seats are divided into 2 train sections Section A and Section B, each is represented by a Boolean array of length 20
// the allocation is done sequentially first from Section A0 to A19 and the B0 to B19 with seat number format as 'A1'
func (s *BookingServer) allocateSeat(train_id string) (seat_number string, err error) {

	for _, train := range *s.trainChart {
		if train.TrainID == train_id {

			for seatID, reserverd := range train.SectionA {
				if !reserverd {
					// setting the seat in Section A to booked
					train.SectionA[seatID] = true
					seat_number := fmt.Sprintf("%s%d", "A", seatID)
					return seat_number, nil
				}
			}

			for seatID, reserverd := range train.SectionB {
				if !reserverd {
					// setting the seat in Section B to booked
					train.SectionB[seatID] = true
					seat_number := fmt.Sprintf("%s%d", "B", seatID)
					return seat_number, nil

				}
			}
		}
	}

	// no available seat found, returning error
	err = errors.New("error allocating seat, no seat available")
	return

}

// deallocateSeats deallocates seats when user requests to cancel for a given train ID and a seat number like "A0".
func (s *BookingServer) deallocateSeat(train_id string, seat_number string) error {
	// get the train section
	section := seat_number[0]
	// get the seat number
	seatID, err := strconv.Atoi(seat_number[1:])

	if err != nil {
		return errors.New("invalid seat ID")
	}
	// Loop through the train routes to find the matching train ID.
	for _, train := range *s.trainChart {
		if train.TrainID == train_id {
			if section == 'A' {
				// mark the seat as available in section A
				train.SectionA[seatID] = false
				return nil
			} else if section == 'B' {
				// mark the seat as available in section B
				train.SectionB[seatID] = false
				return nil
			}
		}
	}
	return errors.New("error de-allocating seat")
}

// modifySeats takes in the old seat number and modifies to the new requested one if it is available
func (s *BookingServer) modifySeat(train_id string, seat_number string, requested_seat string) error {

	// get the old train seat details
	old_section := seat_number[0]
	old_seatID, err := strconv.Atoi(seat_number[1:])
	if err != nil {
		return errors.New("invalid old seat ID")
	}
	// get the new train seat details
	new_section := requested_seat[0]
	new_seatID, err := strconv.Atoi(requested_seat[1:])
	if err != nil {
		return errors.New("invalid new seat ID")
	}

	new_seat_available := false

	// Loop through the train routes to find if the requested seat is available.
	for _, train := range *s.trainChart {
		if train.TrainID == train_id {
			if new_section == 'A' {
				if !train.SectionA[new_seatID] {
					// allot the requested seat
					new_seat_available = true
					train.SectionA[new_seatID] = true
				}
			} else if new_section == 'B' {
				if !train.SectionB[new_seatID] {
					// allot the requested seat
					new_seat_available = true
					train.SectionB[new_seatID] = true
				}
			}
			if new_seat_available {
				if old_section == 'A' {
					// deallocate the old seat
					train.SectionA[old_seatID] = false
					return nil
				} else if old_section == 'B' {
					// deallocate the old seat
					train.SectionB[old_seatID] = false
					return nil
				}
			}

		}
	}
	return errors.New("error modifying seat, requested seat not available")
}

// get trainBySeatSection returns a tainID along with list of users and their seat numbers
func (s *BookingServer) getTrainChart(from string, to string) (train_id string, chart string, err error) {

	for _, train := range *s.trainChart {
		if train.From == from && train.To == to {
			return train.TrainID, createChart(train.SectionA[:], train.SectionB[:]), nil
		}
	}
	// No train found, default return values are used
	err = errors.New("no trains for given route")
	return
}

// createChart returns a reserved seats chart
func createChart(sectionA []bool, sectionB []bool) string {

	var chart1, chart2 strings.Builder

	// Generate chart for the first array
	for i, reserved := range sectionA {
		if reserved {
			chart1.WriteString(fmt.Sprintf("A%d reserved; ", i))
		}
	}

	// Generate chart for the second array
	for i, reserved := range sectionB {
		if reserved {
			chart2.WriteString(fmt.Sprintf("B%d reserved; ", i))
		}
	}

	// Trim the trailing semicolon and space
	chart1Str := strings.TrimSpace(chart1.String())
	if len(chart1Str) > 0 && chart1Str[len(chart1Str)-1:] == ";" {
		chart1Str = chart1Str[:len(chart1Str)-1]
	}

	chart2Str := strings.TrimSpace(chart2.String())
	if len(chart2Str) > 0 && chart2Str[len(chart2Str)-1:] == ";" {
		chart2Str = chart2Str[:len(chart2Str)-1]
	}

	// Concatenate the two charts with a newline
	result := fmt.Sprintf("%s\n%s", chart1Str, chart2Str)
	return result

}

// MockTrainData is mock data of train routes
func MockTrainData() *[]*train.TrainChart {
	return &[]*train.TrainChart{
		{
			TrainID: "Route1",
			From:    "London",
			To:      "France",
			Price:   20.0,
			// SectionA and SectionB are default-initialized to false
		},
		{
			TrainID: "Route2",
			From:    "CityA",
			To:      "CityB",
			Price:   120.0,
		},
		{
			TrainID: "Route3",
			From:    "CityA",
			To:      "CityD",
			Price:   140.0,
		},
		{
			TrainID: "Route4",
			From:    "CityB",
			To:      "CityC",
			Price:   110.0,
		},
		{
			TrainID: "Route5",
			From:    "CityB",
			To:      "CityD",
			Price:   130.0,
		},
		{
			TrainID: "Route6",
			From:    "CityC",
			To:      "CityD",
			Price:   150.0,
		},
		{
			TrainID: "Route7",
			From:    "CityA",
			To:      "CityE",
			Price:   160.0,
		},
		{
			TrainID: "Route8",
			From:    "CityB",
			To:      "CityE",
			Price:   170.0,
		},
		{
			TrainID: "Route9",
			From:    "CityC",
			To:      "CityE",
			Price:   180.0,
		},
		{
			TrainID: "Route10",
			From:    "CityD",
			To:      "CityE",
			Price:   190.0,
		},
	}
}
