package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
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

	// ------------------- Rooms Endpoints -------------------
	router.GET("/rooms/:room_type", getAvailableRoomsByType)

	return router
}

func TestAllEndpoints(t *testing.T) {
	router := setupRouter() // Assuming setupRouter() initializes your routes

	tests := []struct {
		name         string
		method       string
		url          string
		body         interface{}
		expectStatus int
	}{
		{"Get Guests", "GET", "/guests", nil, http.StatusOK},
		{"Get Guest by ID", "GET", "/guests/1", nil, http.StatusOK},
		{"Create Guest", "POST", "/guests", map[string]string{"name": "John Doe", "email": "john@example.com"}, http.StatusCreated},
		{"Update Guest", "PUT", "/guests/1", map[string]string{"name": "John Smith"}, http.StatusOK},
		{"Delete Guest", "DELETE", "/guests/1", nil, http.StatusNoContent},

		{"Get Payments", "GET", "/payments", nil, http.StatusOK},
		{"Get Payment by ID", "GET", "/payments/1", nil, http.StatusOK},
		{"Create Payment", "POST", "/payments", map[string]interface{}{"amount": 100.0, "guest_id": 1}, http.StatusCreated},
		{"Update Payment", "PUT", "/payments/1", map[string]interface{}{"amount": 150.0}, http.StatusOK},
		{"Delete Payment", "DELETE", "/payments/1", nil, http.StatusNoContent},

		{"Get Room by Type", "GET", "/rooms/Single", nil, http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if test.body != nil {
				jsonBody, _ := json.Marshal(test.body)
				req, err = http.NewRequest(test.method, test.url, bytes.NewBuffer(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(test.method, test.url, nil)
			}

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)
			if recorder.Code != test.expectStatus {
				t.Errorf("%s: expected status %d but got %d", test.name, test.expectStatus, recorder.Code)
			}
		})
	}
}
