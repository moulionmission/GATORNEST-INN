// server.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Package struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var packages = []Package{
	{ID: 1, Name: "Deluxe Room", Price: 100},
	{ID: 2, Name: "Suite", Price: 200},
}

func getPackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}

func getPackageByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/api/packages/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for _, pkg := range packages {
		if pkg.ID == id {
			json.NewEncoder(w).Encode(pkg)
			return
		}
	}
	http.Error(w, "Package not found", http.StatusNotFound)
}

func createPackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newPkg Package
	if err := json.NewDecoder(r.Body).Decode(&newPkg); err != nil || newPkg.Name == "" || newPkg.Price == 0 {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	newPkg.ID = len(packages) + 1
	packages = append(packages, newPkg)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPkg)
}

func main() {
	http.HandleFunc("/api/packages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPackages(w, r)
		} else if r.Method == http.MethodPost {
			createPackage(w, r)
		}
	})
	http.HandleFunc("/api/packages/", getPackageByID)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
