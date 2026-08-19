# Train Ticketing System

This project is a simple train ticket booking system built using Go. It allows users to book train tickets by communicating with a gRPC-based server.

## Features

- **Ticket Booking**: Users can book tickets for a train route from a predefined list of route.
- **gRPC Communication**: The client and server communicate using gRPC, providing a robust and scalable way to handle ticket booking requests.
- **Show Receipt**: User can see the receipt of their booked ticket.
- **Modify ticket**: User can request to modify the seat allotted to them.
- **Remove Ticket**: Use can request to cancel the ticket.
- **View Ticket by train Section**: User can view the allotted seats of a train by section and username.
- **Concurrency Control**: The system can enable concurrency that allows only one user to book a ticket at a time.


## How to Run the Project

### Prerequisites

- Go
- Git

### Steps to Run

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/saumeya/train-ticketing.git
   cd train-ticketing
2. **Run the Server**:
   ```bash
   go run server/main.go
3. **Run the Client on a different Terminal**:
   ```bash
   go run client/main.go
   
## Assumptions

- **Single Ticket Booking**: The User can book only one ticket at a time.
- **Predefined Train Routes**: The system uses a predefined list of train routes stored in a mock database. The routes are initialized in the code, and users can only book tickets for these routes.
- **Seat Allotment Process**: Seat Allotment is sequentially implemented, the next available seat is given to the user starting from Section A to Section B
- **Email as Username**: Each user is identified by a unique email address, which serves as their username for all train seat operations.
- **Train Seats and Sections**: Each train is divided into two sections, A and B, with 20 seats available in each section.

## Components

- **Server**: Handles incoming ticket booking requests from clients and processes them. It includes separation between booking services and train operations and unit tests.
- **Client**: Sends ticket booking requests to the server and displays the server's response.
- **API**: Contains the protobuf (proto) definitions and structs used for gRPC communication. This includes message types, service definitions, and RPC methods.
  
## Future Scope

- **User Authentication**: Implement user authentication to allow multiple users to securely book tickets. Add login/signup options.
- **Database Integration**: Replace the mock database with a real database (e.g., PostgreSQL, MySQL) to store user information and ticket bookings persistently.
- **Multiple Ticket Booking**: Allow users to book multiple tickets in a single transaction, improving flexibility and user experience.
- **Improved Seat Allotement**: Develop a better algorithm for seat allocation that considers user preferences and introduces randomization to optimize seat assignments.

## Sample Output

```bash
------------------------------- 

Book Ticket Service Call 
Booking ID: f74d53e7-bfe3-4587-93fe-dc2a048a921d
From: London
To: France
First Name: Sam
Last Name: Hill
Email: samhill@gmail.com
Price Paid: 20.00
Seat Number: A0

------------------------------- 

Book Ticket Service Call 
Booking ID: 1f2c5bb9-df54-4244-825a-d7c24a42bbef
From: London
To: France
First Name: Tom
Last Name: Holl
Email: tomholl@gmail.com
Price Paid: 20.00
Seat Number: A1

------------------------------- 

View Train Chart 
From: London
To: France
TrainID: Route1
Chart: 
A0 reserved; A1 reserved


------------------------------- 

Show receipt Service Call 
Booking ID: f74d53e7-bfe3-4587-93fe-dc2a048a921d
From: London
To: France
First Name: Sam
Last Name: Hill
Email: samhill@gmail.com
Price Paid: 20.00
Seat Number: A0

------------------------------- 

View Seats by Section Service Call 
TrainID: Route1
From: London
To: France
Train Section: SectionA
User Seat Mapping:
Seat Number: A0, FirstName: Sam, LastName: Hill, Email: samhill@gmail.com
Seat Number: A1, FirstName: Tom, LastName: Holl, Email: tomholl@gmail.com

------------------------------- 

Modify Ticket Service Call 
2024/09/03 11:02:39 Seat is successfully updated for Booking ID 12345454 
2024/09/03 11:02:39 Old Seat: A0
New Seat: A3

------------------------------- 

Remove user ticket Service Call 
2024/09/03 11:02:39 Seat is successfully cancelled for Booking ID 12345454 
2024/09/03 11:02:39 Cancelled Seat: A3

Attempting to run thread ID 3 
Booking ID: 09476702-33e6-47ab-8ec4-b6c5ab3c4507
From: London
To: France
First Name: sam
Last Name: doen
Email: samdoen@gmail.com
Price Paid: 20.00
Seat Number: A0

Attempting to run thread ID 2 
Booking ID: ec66c339-3732-48a7-a624-ddf8979daae1
From: London
To: France
First Name: sam
Last Name: doen
Email: samdoen@gmail.com
Price Paid: 20.00
Seat Number: A2

Attempting to run thread ID 1 
Booking ID: fee1c28e-479e-4678-9236-85e208e80f81
From: London
To: France
First Name: sam
Last Name: doen
Email: samdoen@gmail.com
Price Paid: 20.00
Seat Number: A3
