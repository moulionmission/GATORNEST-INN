package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

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
	ReservationID int       `json:"reservation_id"`
	GuestID       int       `json:"guest_id"`
	RoomID        int       `json:"room_id"`
	CheckInDate   string    `json:"check_in_date"`
	CheckOutDate  string    `json:"check_out_date"`
	Status        string    `json:"status"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        int       `json:"user_id,omitempty"` // added for linking to logged-in user
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
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenStr := authHeader
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
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

	result, err := db.Exec(`
		INSERT INTO Reservations (guest_id, room_id, check_in_date, check_out_date, status, total_price, user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		reservation.GuestID, reservation.RoomID, reservation.CheckInDate,
		reservation.CheckOutDate, reservation.Status, reservation.TotalPrice, reservation.UserID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating reservation", "details": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	reservation.ReservationID = int(id)

	c.JSON(http.StatusCreated, reservation)
}
