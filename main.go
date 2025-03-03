package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

// Global DB variable
var db *sql.DB

// Guest represents the guest profile as defined in the Guests table.
type Guest struct {
	GuestID   int    `json:"guest_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
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

	// ------------------- Guest Endpoints -------------------
	router.POST("/guests", createGuest)
	router.GET("/guests", getGuests)
	router.GET("/guests/:id", getGuest)
	router.PUT("/guests/:id", updateGuest)
	router.DELETE("/guests/:id", deleteGuest)

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

	result, err := db.Exec("INSERT INTO Guests (first_name, last_name, email, phone) VALUES (?, ?, ?, ?)",
		guest.FirstName, guest.LastName, guest.Email, guest.Phone)
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
