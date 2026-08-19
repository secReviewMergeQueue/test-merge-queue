package main

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/saumeya/train-ticketing/api/proto"

	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var mu sync.Mutex

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	defer conn.Close()
	client := pb.NewBookingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	bookTicket(ctx, client, &pb.User{FirstName: "Sam", LastName: "Hill", Email: "samhill@gmail.com"})
	bookTicket(ctx, client, &pb.User{FirstName: "Tom", LastName: "Holl", Email: "tomholl@gmail.com"})
	viewTrainChart(ctx, client)
	showReceipt(ctx, client)
	viewSeatsBySection(ctx, client)
	modifyUserSeat(ctx, client)
	removeUser(ctx, client)

	// adding concurrency for ticket booking
	var wg sync.WaitGroup

	// Launch multiple goroutines to book tickets concurrently
	for i := 1; i <= 3; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// Create a context with a timeout
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			bookTicketConcurrently(timeoutCtx, client, id)
		}(i)
	}
	// Wait for all the goroutines to complete before exiting
	wg.Wait()
}

// bookTikcet calls the ticket booking service
func bookTicket(ctx context.Context, client pb.BookingServiceClient, user *pb.User) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nBook Ticket Service Call ")

	response, err := client.BookTicket(ctx, &pb.BookingRequest{
		From: "London",
		To:   "France",
		User: user,
	})
	if err != nil {
		log.Fatalf("Error booking the ticket: %v", err)
	}

	printBookingResponse(response)
}

// showReceipt call the show receipt service
func showReceipt(ctx context.Context, client pb.BookingServiceClient) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nShow receipt Service Call ")

	response, err := client.ShowReceipt(ctx, &pb.ShowReceiptRequest{
		BookingId: "12345454",
		User: &pb.User{
			FirstName: "sam",
			LastName:  "hill",
			Email:     "samhill@gmail.com",
		},
	})
	if err != nil {
		log.Fatalf("Error getting your booking receipt: %v \n", err)
	}

	printBookingResponse(response)
}

// removeUser calls the remove user or cancel ticket service
func removeUser(ctx context.Context, client pb.BookingServiceClient) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nRemove user ticket Service Call ")

	r, err := client.RemoveUser(ctx, &pb.RemoveUserRequest{
		BookingId: "12345454",
		User: &pb.User{
			FirstName: "sam",
			LastName:  "hill",
			Email:     "samhill@gmail.com",
		},
	})
	if err != nil {
		log.Fatalf("Error cancelling the ticket: %v", err)
	}

	log.Printf("Seat is successfully cancelled for Booking ID %s \n", r.GetBookingId())
	log.Printf("Cancelled Seat: %s\n", r.GetSeatNumber())
}

// modifyUser calls the modify user seat service
func modifyUserSeat(ctx context.Context, client pb.BookingServiceClient) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nModify Ticket Service Call ")

	r, err := client.ModifyUserSeat(ctx, &pb.ModifyUserSeatRequest{
		BookingId: "12345454",
		User: &pb.User{
			FirstName: "sam",
			LastName:  "hill",
			Email:     "samhill@gmail.com",
		},
		RequestedSeat: "A3",
	})
	if err != nil {
		log.Fatalf("Error modifying the ticket: %v", err)
	}

	log.Printf("Seat is successfully updated for Booking ID %s \n", r.GetBookingId())
	log.Printf("Old Seat: %s\nNew Seat: %s\n", r.GetSeatNumber(), r.GetRequestSeat())
}

// viewSeatBySection calls the view seat by section service
func viewSeatsBySection(ctx context.Context, client pb.BookingServiceClient) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nView Seats by Section Service Call ")

	r, err := client.ViewSeatsBySection(ctx, &pb.ViewSeatsBySectionRequest{
		From:         "London",
		To:           "France",
		TrainSection: "SectionA",
	})
	if err != nil {
		log.Fatalf("error fetching the train data: %v", err)
	}

	printBookingResponse(r)
}

// viewTrainChart calls the view train chart  service
func viewTrainChart(ctx context.Context, client pb.BookingServiceClient) {

	fmt.Println("\n------------------------------- ")
	fmt.Println("\nView Train Chart ")

	r, err := client.ViewTrainChart(ctx, &pb.ViewTrainChartRequest{
		From: "London",
		To:   "France",
	})
	if err != nil {
		log.Fatalf("error fetching the train data: %v", err)
	}

	printBookingResponse(r)
}

func bookTicketConcurrently(ctx context.Context, client pb.BookingServiceClient, id int) {
	// Acquire the mutex lock
	mu.Lock()

	// Unlock when function exits
	defer mu.Unlock()

	fmt.Printf("Attempting to run thread ID %d \n", id)

	select {
	case <-time.After(2 * time.Second): // Simulate booking process time
		response, err := client.BookTicket(ctx, &pb.BookingRequest{
			From: "London",
			To:   "France",
			User: &pb.User{
				FirstName: "sam",
				LastName:  "doen",
				Email:     "samdoen@gmail.com",
			},
		})
		if err != nil {
			log.Fatalf("error booking the ticket: %v", err)
		}

		printBookingResponse(response)
	case <-ctx.Done(): // If context is canceled due to timeout

	}

}

func printBookingResponse(r interface{}) {
	switch response := r.(type) {
	case *pb.BookingResponse:
		fmt.Println("Booking ID: " + response.GetBookingId())
		fmt.Println("From: " + response.GetFrom())
		fmt.Println("To: " + response.GetTo())
		fmt.Println("First Name: " + response.GetUser().GetFirstName())
		fmt.Println("Last Name: " + response.GetUser().GetLastName())
		fmt.Println("Email: " + response.GetUser().GetEmail())
		fmt.Printf("Price Paid: %.2f\n", response.GetPricePaid())
		fmt.Println("Seat Number: " + response.GetSeatNumber())
	case *pb.ShowReceiptResponse:
		fmt.Println("Booking ID: " + response.GetBookingId())
		fmt.Println("From: " + response.GetFrom())
		fmt.Println("To: " + response.GetTo())
		fmt.Println("First Name: " + response.GetUser().GetFirstName())
		fmt.Println("Last Name: " + response.GetUser().GetLastName())
		fmt.Println("Email: " + response.GetUser().GetEmail())
		fmt.Printf("Price Paid: %.2f\n", response.GetPricePaid())
		fmt.Println("Seat Number: " + response.GetSeatNumber())
	case *pb.ViewTrainChartResponse:
		fmt.Println("From: " + response.GetFrom())
		fmt.Println("To: " + response.GetTo())
		fmt.Println("TrainID: " + response.GetTrainId())
		fmt.Println("Chart: \n" + response.GetResponse())
	case *pb.ViewSeatsBySectionResponse:
		fmt.Println("TrainID: " + response.GetTrainId())
		fmt.Println("From: " + response.GetFrom())
		fmt.Println("To: " + response.GetTo())
		fmt.Println("Train Section: " + response.GetTrainSection())
		fmt.Println("User Seat Mapping:")
		for _, seat := range response.GetSeats() {
			user := seat.GetUser()
			seatInfo := fmt.Sprintf(
				"Seat Number: %s, FirstName: %s, LastName: %s, Email: %s",
				seat.GetSeatNumber(),
				user.GetFirstName(),
				user.GetLastName(),
				user.GetEmail(),
			)
			fmt.Println(seatInfo)
		}

	default:
		fmt.Println("Unknown Type")
	}

}
