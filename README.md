Bus Reservation System (Demo Project)

  This Bus Reservation System was created as a demo project to showcase my Golang skills, particularly focusing on REST API development and database functionality. I built 
this project within a short time frame after realizing that Dubai recruiters often look for GitHub projects in a candidate's portfolio.

Project Overview
Implements REST APIs for bus booking and management.
Uses MySQL as the database to store user details, bus information, and booking records.
Designed with a layered architecture for scalability and maintainability.
Effectively handles all operations, ensuring smooth data flow between APIs, business logic, and the database.

Project Structural layers:
1.  db - Handles database operations and queries.
2. bl - Implements core business logic for the system.
3. endpoints - Defines REST API endpoints for user interactions.
4. models - Represents data models and database mappings.

Below are a few exposed REST APIS:
1. Register a User (POST /register)
2. Get All Buses (GET /buses)      
3. Add a New Bus (POST /bus)     
4. Check Available Seats for a Bus (GET /availableSeats/{bus_id})
5. Returns the number of unreserved seats for a specific bus.
6. Get Total Bookings for a Bus (GET /totalBookings/{bus_id})
7. Fetches the total number of bookings made for a particular bus.
8. Book Seats on a Bus (POST /bookseats) 

To run this project you just need MySQL already installed. 
This project is a technical demonstration of my ability to structure and implement backend systems efficiently.
