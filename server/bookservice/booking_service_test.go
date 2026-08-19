package bookingservice

import (
	"context"
	"testing"

	pb "github.com/saumeya/train-ticketing/api/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBookTicket_InvalidRequest(t *testing.T) {
	server := NewBookingService()
	request := &pb.BookingRequest{
		From: "", // Invalid from
		To:   "CityB",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	response, err := server.BookTicket(context.Background(), request)

	assert.Nil(t, response)
	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestBookTicket_Success(t *testing.T) {
	server := NewBookingService()
	request := &pb.BookingRequest{
		From: "CityA",
		To:   "CityB",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	response, err := server.BookTicket(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "CityA", response.GetFrom())
	assert.Equal(t, "CityB", response.GetTo())
	assert.Equal(t, "John", response.GetUser().GetFirstName())
	assert.NotEmpty(t, response.GetBookingId())
	assert.NotEmpty(t, response.GetSeatNumber())
	assert.Greater(t, response.GetPricePaid(), 0.0)
}

func TestShowReceipt_Success(t *testing.T) {
	server := NewBookingService()

	// Book a ticket first to generate a reservation
	request := &pb.BookingRequest{
		From: "CityA",
		To:   "CityB",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}
	response, err := server.BookTicket(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	// Now show the receipt
	receiptRequest := &pb.ShowReceiptRequest{
		BookingId: "123536",
		User: &pb.User{
			Email: "john.doe@example.com",
		},
	}
	receipt, err := server.ShowReceipt(context.Background(), receiptRequest)

	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, "CityA", receipt.GetFrom())
	assert.Equal(t, "CityB", receipt.GetTo())
	assert.Equal(t, "John", receipt.GetUser().GetFirstName())
}

func TestRemoveUser_Success(t *testing.T) {
	server := NewBookingService()

	// Book a ticket first to generate a reservation
	request := &pb.BookingRequest{
		From: "CityA",
		To:   "CityB",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}
	server.BookTicket(context.Background(), request)

	// Now remove the user
	removeRequest := &pb.RemoveUserRequest{
		BookingId: "123536",
		User: &pb.User{
			Email: "john.doe@example.com",
		},
	}
	removeResponse, err := server.RemoveUser(context.Background(), removeRequest)

	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)
}

func TestModifyUserSeat_Success(t *testing.T) {
	server := NewBookingService()

	// Book a ticket first to generate a reservation
	request := &pb.BookingRequest{
		From: "CityA",
		To:   "CityB",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}
	server.BookTicket(context.Background(), request)

	// Now modify the seat
	modifyRequest := &pb.ModifyUserSeatRequest{
		BookingId: "123536",
		User: &pb.User{
			Email: "john.doe@example.com",
		},
		RequestedSeat: "B1", // Assuming "B1" is a valid and available seat
	}
	modifyResponse, err := server.ModifyUserSeat(context.Background(), modifyRequest)

	assert.NoError(t, err)
	assert.NotNil(t, modifyResponse)
}
