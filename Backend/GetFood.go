package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Food struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name"`
	Description1 string `json:"description1"`
	Description2 string `json:"description2"`
	Description3 string `json:"description3"`
	Description4 string `json:"description4"`
	ImageURL     string `json:"image_url"`
	Country      string `json:"country"`
}

var db *gorm.DB

func init() {
	var err error

	// Подключение к базе данных
	dsn := "postgres://postgres.omqkkeruydkttwwkdnib:50TADocqYFe4CFTx@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?sslmode=require&supa=base-pooler.x"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Автоматическая миграция для структуры Food
	if err := db.AutoMigrate(&Food{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}

func getFoodByName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodName := params["name"]
	if foodName == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Name parameter is required")
		return
	}

	fmt.Printf("Searching for food: %s\n", foodName)

	var food Food
	if err := db.Where("name = ?", foodName).First(&food).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Couldn't find food: %s\n", foodName)
			sendErrorResponse(w, http.StatusNotFound, "Food not found")
		} else {
			fmt.Printf("Database error while searching for food: %s\n", foodName)
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error querying the database: %v", err))
		}
		return
	}

	fmt.Printf("Found food: %s (ID: %d)\n", food.Name, food.ID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(food); err != nil {
		fmt.Printf("Error encoding JSON for food: %s\n", foodName)
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error encoding JSON: %v", err))
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(errorResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/food/{name}", getFoodByName).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(origins, methods, headers)(r)); err != nil {
		log.Fatal("Server failed: ", err)
	}
}
