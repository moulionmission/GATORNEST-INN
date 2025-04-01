# 🏨 Gatorrest Inn – Backend API

This is the backend system for **Gatorrest Inn**, a hotel management application.  
Built using **Go (Gin framework)** and **MySQL**, this phase focuses on secure user authentication and linking guests and reservations to logged-in users.

---

## 🔧 Features

- ✅ User registration with **email + password**
- ✅ **JWT authentication**
- ✅ Protected API routes
- ✅ Guests linked to users
- ✅ Reservations linked to users
- ✅ `/profile` endpoint to fetch current user info

---

## 🚀 Tech Stack

| Component     | Tech                 |
|---------------|----------------------|
| Language      | Go (Golang)          |
| Framework     | Gin                  |
| Database      | MySQL                |
| Auth          | JWT (golang-jwt)     |
| ORM/SQL       | database/sql + raw   |
| Environment   | godotenv             |
| API Testing   | Postman              |

---

## 📁 Project Structure

. ├── main.go # Main API logic ├── go.mod / go.sum # Dependency management ├── hotel_db_setup.sql # Database schema ├── .env.example # Environment variable template


🔐 API Endpoints
🔸 Authentication
Method	Endpoint	Description
POST	/register	Register a new user
POST	/login	Log in, get JWT
GET	/profile	Get current user info
Requires Authorization header with JWT for protected routes.

🔸 Guests (Protected)
Method	Endpoint	Description
POST	/guests	Create a new guest
GET	/guests	List all your guests
GET	/guests/:id	Get a guest by ID
PUT	/guests/:id	Update guest info
DELETE	/guests/:id	Delete a guest
🔸 Reservations (Protected)
Method	Endpoint	Description
POST	/reservations	Create a new reservation
(GET)	(coming soon)	Filter by logged-in user
🧠 Project Phase
This backend was developed as part of the Gatorrest Inn project, with the objective to:

Build a secure, token-based user system

Link data to authenticated users

Provide a functional backend ready for frontend integration.
