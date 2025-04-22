package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Global DB variable
var db *sql.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Guest represents the guest profile as defined in the Guests table.
type Guest struct {
	GuestID   int    `json:"guest_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserID    int    `json:"user_id,omitempty"` // Set automatically from JWT
}

type User struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"` // Raw password from request (not stored)
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Payment represents the payment record as defined in the Payments table.
type Payment struct {
	PaymentID       int       `json:"payment_id"`
	ReservationID   int       `json:"reservation_id"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentStatus   string    `json:"payment_status"`
	Amount          float64   `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
}

type Reservation struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ReservationID int       `json:"reservation_id"`
	GuestID       int       `json:"guest_id"`
	RoomID        int       `json:"room_id"`
	CheckInDate   string    `json:"check_in_date"`
	CheckOutDate  string    `json:"check_out_date"`
	Status        string    `json:"status"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        int       `json:"user_id,omitempty"` // added for linking to logged-in user
	Email         string    `json:"email"`
}

type Schedule struct {
	ScheduleID int    `json:"schedule_id,omitempty"`
	StaffID    int    `json:"staff_id"`
	ShiftDate  string `json:"shift_date"` // "Monday", "Tuesday", etc.
	ShiftTime  string `json:"shift_time"` // "Morning", "Afternoon", "Night"
}

// initDB loads environment variables and connects to the MySQL database.
func initDB() {
	// Load environment variables from .env (if present)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, dbname)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging DB:", err)
	}
	log.Println("Connected to the database successfully")
}

func main() {
	initDB()
	router := gin.Default()

	// Apply CORS middleware (allowing requests from http://localhost:3001)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"}, // React frontend's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.POST("/register", registerUser)
	router.POST("/login", loginUser)

	// Protected Guest routes
	auth := router.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/guests", createGuest)
		auth.GET("/guests", getGuests)
		auth.GET("/guests/:id", getGuest)
		auth.PUT("/guests/:id", updateGuest)
		auth.DELETE("/guests/:id", deleteGuest)

		auth.GET("/profile", getProfile)
		auth.POST("/reservations", createReservation)

		auth.GET("/staff/:id", getStaffByID)
		auth.GET("/staffs", getAllStaff)

		auth.POST("/schedule", addToSchedule)
		auth.DELETE("/schedule/:id", removeFromSchedule)
		auth.PUT("/schedule/:id", updateSchedule)
		auth.GET("/schedules", getAllSchedules)

	}

	// ------------------- Payment Endpoints -------------------
	router.POST("/payments", createPayment)
	router.GET("/payments", getPayments)
	router.GET("/payments/:id", getPayment)
	router.PUT("/payments/:id", updatePayment)
	router.DELETE("/payments/:id", deletePayment)

	// Start the server on PORT defined in .env (default 3000)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}

// GenerateRandomPhoneNumber creates a random 10-digit phone number
func GenerateRandomPhoneNumber() string {
	rand.Seed(time.Now().UnixNano())

	// Generate random digits
	areaCode := rand.Intn(900) + 100 // Ensures a 3-digit number (100-999)
	prefix := rand.Intn(900) + 100   // Ensures a 3-digit number (100-999)
	lineNumber := rand.Intn(10000)   // Ensures a 4-digit number (0000-9999)

	// Format as a phone number
	return fmt.Sprintf("%d-%d-%04d", areaCode, prefix, lineNumber)
}

// ------------------- Guest Handlers -------------------

// createGuest handles POST /guests
func createGuest(c *gin.Context) {
	var guest Guest
	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// Get user_id from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := userIDInterface.(int)
	guest.UserID = userID

	result, err := db.Exec(
		"INSERT INTO Guests (first_name, last_name, email, phone, user_id) VALUES (?, ?, ?, ?, ?)",
		guest.FirstName, guest.LastName, guest.Email, guest.Phone, guest.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating guest", "details": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving guest id", "details": err.Error()})
		return
	}
	guest.GuestID = int(id)

	c.JSON(http.StatusCreated, guest)
}

// getProfile handles GET /profile
func getProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}

// getGuests handles GET /guests
func getGuests(c *gin.Context) {
	rows, err := db.Query("SELECT guest_id, first_name, last_name, email, phone FROM Guests")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching guests", "details": err.Error()})
		return
	}
	defer rows.Close()

	var guests []Guest
	for rows.Next() {
		var g Guest
		if err := rows.Scan(&g.GuestID, &g.FirstName, &g.LastName, &g.Email, &g.Phone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning guest", "details": err.Error()})
			return
		}
		guests = append(guests, g)
	}
	c.JSON(http.StatusOK, guests)
}

// getGuest handles GET /guests/:id
func getGuest(c *gin.Context) {
	id := c.Param("id")
	var g Guest
	err := db.QueryRow("SELECT guest_id, first_name, last_name, email, phone FROM Guests WHERE guest_id = ?", id).
		Scan(&g.GuestID, &g.FirstName, &g.LastName, &g.Email, &g.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching guest", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, g)
}

// updateGuest handles PUT /guests/:id
func updateGuest(c *gin.Context) {
	id := c.Param("id")
	var guest Guest
	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	result, err := db.Exec("UPDATE Guests SET first_name = ?, last_name = ?, email = ?, phone = ? WHERE guest_id = ?",
		guest.FirstName, guest.LastName, guest.Email, guest.Phone, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating guest", "details": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
		return
	}
	guest.GuestID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, guest)
}

// deleteGuest handles DELETE /guests/:id
func deleteGuest(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM Guests WHERE guest_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting guest", "details": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Guest not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Guest deleted successfully"})
}

// ------------------- Payment Handlers -------------------

// createPayment handles POST /payments
func createPayment(c *gin.Context) {
	var payment Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	result, err := db.Exec("INSERT INTO Payments (reservation_id, payment_method, payment_status, amount) VALUES (?, ?, ?, ?)",
		payment.ReservationID, payment.PaymentMethod, payment.PaymentStatus, payment.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating payment", "details": err.Error()})
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving payment id", "details": err.Error()})
		return
	}
	payment.PaymentID = int(id)
	// The transaction_date is automatically set by the DB.
	c.JSON(http.StatusCreated, payment)
}

// getPayments handles GET /payments
func getPayments(c *gin.Context) {
	rows, err := db.Query("SELECT payment_id, reservation_id, payment_method, payment_status, amount, transaction_date FROM Payments")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching payments", "details": err.Error()})
		return
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.PaymentID, &p.ReservationID, &p.PaymentMethod, &p.PaymentStatus, &p.Amount, &p.TransactionDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning payment", "details": err.Error()})
			return
		}
		payments = append(payments, p)
	}
	c.JSON(http.StatusOK, payments)
}

// getPayment handles GET /payments/:id
func getPayment(c *gin.Context) {
	id := c.Param("id")
	var p Payment
	err := db.QueryRow("SELECT payment_id, reservation_id, payment_method, payment_status, amount, transaction_date FROM Payments WHERE payment_id = ?", id).
		Scan(&p.PaymentID, &p.ReservationID, &p.PaymentMethod, &p.PaymentStatus, &p.Amount, &p.TransactionDate)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching payment", "details": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, p)
}

// updatePayment handles PUT /payments/:id
func updatePayment(c *gin.Context) {
	id := c.Param("id")
	var payment Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	result, err := db.Exec("UPDATE Payments SET reservation_id = ?, payment_method = ?, payment_status = ?, amount = ? WHERE payment_id = ?",
		payment.ReservationID, payment.PaymentMethod, payment.PaymentStatus, payment.Amount, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating payment", "details": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	payment.PaymentID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, payment)
}

// deletePayment handles DELETE /payments/:id
func deletePayment(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM Payments WHERE payment_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting payment", "details": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}

// isValidEmail performs basic regex check for valid email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// registerUser handles POST /register
func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	result, err := db.Exec("INSERT INTO Users (email, password_hash) VALUES (?, ?)", user.Email, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user", "details": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	user.UserID = int(id)
	user.Password = "" // Clear password before sending back

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.UserID,
	})
}

// loginUser handles POST /login
func loginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if !isValidEmail(loginData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var storedHash string
	var userID int

	err := db.QueryRow("SELECT user_id, password_hash FROM Users WHERE email = ?", loginData.Email).
		Scan(&userID, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during login", "details": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   loginData.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// AuthMiddleware verifies the JWT token and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Received Authorization Header:", authHeader) // Debugging line

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Println("Extracted Token:", tokenString) // Debugging line

		// Debugging: Decode the token before validation
		claims := jwt.MapClaims{}
		_, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims)
		if err != nil {
			fmt.Println("Error parsing token:", err)
		} else {
			fmt.Println("Decoded Token Claims:", claims)
		}

		// tokenStr := authHeader
		// token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method")
		// 	}
		// 	return jwtSecret, nil
		// })
		// Proceed with actual token validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		// if err != nil || !token.Valid {
		// 	fmt.Println("Token Validation Error:", err) // Debugging line
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		// 	c.Abort()
		// 	return
		// }

		if err != nil {
			fmt.Println("Token Validation Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		// 	c.Abort()
		// 	return
		// }

		// Extract user ID from token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := int(claims["user_id"].(float64)) // Ensure the claim exists
			c.Set("user_id", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("email", claims["email"].(string))

		c.Next()
	}
}

func createReservation(c *gin.Context) {
	var reservation Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// Extract user_id from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	reservation.UserID = userIDInterface.(int)

	// 🔹 Step 1: Check if guest exists
	var existingGuest Guest
	err := db.QueryRow("SELECT guest_id FROM Guests WHERE guest_id = ?", reservation.GuestID).Scan(&existingGuest.GuestID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Guest does not exist, create a new one
			fmt.Println("Guest not found. Creating new guest...")

			newGuest := Guest{
				FirstName: reservation.FirstName,
				LastName:  reservation.LastName,
				Email:     reservation.Email,
				Phone:     GenerateRandomPhoneNumber(),
				UserID:    reservation.UserID,
			}

			result, err := db.Exec(`
				INSERT INTO Guests (first_name, last_name, email, phone, user_id) 
				VALUES (?, ?, ?, ?, ?)`,
				newGuest.FirstName, newGuest.LastName, newGuest.Email, newGuest.Phone, newGuest.UserID,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating guest", "details": err.Error()})
				return
			}

			// Retrieve the new guest_id
			guestID, _ := result.LastInsertId()
			reservation.GuestID = int(guestID)
		} else {
			// Some other error occurred
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
	}

	// 🔹 Step 2: Insert the reservation
	result, err := db.Exec(`
		INSERT INTO Reservations (guest_id, room_id, check_in_date, check_out_date, status, total_price)
    	VALUES (?, ?, ?, ?, ?, ?)`,
		reservation.GuestID, reservation.RoomID, reservation.CheckInDate,
		reservation.CheckOutDate, reservation.Status, reservation.TotalPrice,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation", "details": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	reservation.ReservationID = int(id)

	c.JSON(http.StatusCreated, reservation)
}

// Staff Endpoints

func getAllStaff(c *gin.Context) {
	rows, err := db.Query(`SELECT staff_id, first_name FROM Staff WHERE role='Housekeeping'`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve staff", "details": err.Error()})
		return
	}
	defer rows.Close()

	var staffList []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName string

		if err := rows.Scan(&id, &firstName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse staff", "details": err.Error()})
			return
		}

		staffList = append(staffList, gin.H{
			"staff_id":   id,
			"first_name":  firstName,
		})
	}

	c.JSON(http.StatusOK, staffList)
}



// GET /staff/:id - Get staff information by staff ID
func getStaffByID(c *gin.Context) {
	id := c.Param("id")

	// Query to get staff details
	query := `SELECT staff_id, first_name, last_name, email, role FROM Staff WHERE staff_id = ?`
	row := db.QueryRow(query, id)

	var staffID int
	var firstName, lastName, email, role string

	// Scan the row into variables
	err := row.Scan(&staffID, &firstName, &lastName, &email, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve staff information", "details": err.Error()})
		}
		return
	}

	// Return the staff information as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"staff_id":   staffID,
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"role":       role,
	})
}

// Staff Reservation Endpoints

// GET /schedules - Get all staff schedules
func getAllSchedules(c *gin.Context) {
	rows, err := db.Query(`SELECT staff_id, shift_date, shift_time FROM Staff_Schedule`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve schedules", "details": err.Error()})
		return
	}
	defer rows.Close()

	var schedules []map[string]interface{}
	for rows.Next() {
		var staffID int
		var shiftDate, shiftTime string

		if err := rows.Scan(&staffID, &shiftDate, &shiftTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse schedule", "details": err.Error()})
			return
		}

		schedules = append(schedules, gin.H{
			"staff_id":   staffID,
			"shift_date": shiftDate,
			"shift_time": shiftTime,
		})
	}

	c.JSON(http.StatusOK, schedules)
}

// POST /schedule - Add staff to schedule
func addToSchedule(c *gin.Context) {
	var s Schedule
	if err := c.ShouldBindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `INSERT INTO Staff_Schedule (staff_id, shift_date, shift_time) VALUES (?, ?, ?)`
	result, err := db.Exec(query, s.StaffID, s.ShiftDate, s.ShiftTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add schedule", "details": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	s.ScheduleID = int(id)
	c.JSON(http.StatusCreated, s)
}

// DELETE /schedule/:id
// Removes a schedule entry for a specific staff member based on staff_id (in URL),
// shift_date and shift_time (as query parameters).
// Example: DELETE /schedule/3?shift_date=Friday&shift_time=Night

func removeFromSchedule(c *gin.Context) {
	staffID := c.Param("id")
	shiftDate := c.Query("shift_date")
	shiftTime := c.Query("shift_time")

	if shiftDate == "" || shiftTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing shift_date or shift_time query parameter"})
		return
	}

	query := `
		DELETE FROM Staff_Schedule
		WHERE staff_id = ? AND shift_date = ? AND shift_time = ?
	`
	res, err := db.Exec(query, staffID, shiftDate, shiftTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete schedule", "details": err.Error()})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule removed"})
}

// PUT /schedule/:id
// Updates the shift_time for a specific staff member's schedule
// using staff_id (in URL) and shift_date (in query parameter).
// Example: PUT /schedule/3?shift_date=Friday&new_shift_time=Afternoon

func updateSchedule(c *gin.Context) {
	id := c.Param("id") // staff_id from URL

	shiftDate := c.Query("shift_date")
	newShiftTime := c.Query("new_shift_time")

	if shiftDate == "" || newShiftTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing shift_date or new_shift_time in query parameters"})
		return
	}

	query := `
		UPDATE Staff_Schedule
		SET shift_time = ?
		WHERE staff_id = ? AND shift_date = ?
	`

	res, err := db.Exec(query, newShiftTime, id, shiftDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule", "details": err.Error()})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shift time updated"})
}