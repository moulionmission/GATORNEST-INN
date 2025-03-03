// server_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPackages(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/packages", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPackages)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	var pkgs []Package
	if err := json.NewDecoder(rr.Body).Decode(&pkgs); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	if len(pkgs) == 0 {
		t.Error("Expected at least one package, got none")
	}
}

func TestGetPackageByIDValid(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/packages/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPackageByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, status)
	}

	var pkg Package
	if err := json.NewDecoder(rr.Body).Decode(&pkg); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	if pkg.Name != "Deluxe Room" {
		t.Errorf("Expected package name 'Deluxe Room', got '%s'", pkg.Name)
	}
}

func TestGetPackageByIDInvalid(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/packages/99", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getPackageByID)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %v for non-existent package, got %v", http.StatusNotFound, status)
	}
}

func TestCreatePackageValid(t *testing.T) {
	newPkg := Package{Name: "Luxury Villa", Price: 500}
	payload, _ := json.Marshal(newPkg)
	req, err := http.NewRequest(http.MethodPost, "/api/packages", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createPackage)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %v, got %v", http.StatusCreated, status)
	}

	var createdPkg Package
	if err := json.NewDecoder(rr.Body).Decode(&createdPkg); err != nil {
		t.Errorf("Error decoding response: %v", err)
	}
	if createdPkg.Name != "Luxury Villa" {
		t.Errorf("Expected package name 'Luxury Villa', got '%s'", createdPkg.Name)
	}
	// Assuming the new ID is the length of the packages slice after insertion.
	if createdPkg.ID != len(packages) {
		t.Errorf("Unexpected package ID: %v", createdPkg.ID)
	}
}

func TestCreatePackageInvalid(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/api/packages", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createPackage)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %v for invalid data, got %v", http.StatusBadRequest, status)
	}
}
