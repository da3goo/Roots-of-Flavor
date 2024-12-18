package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"

	_ "github.com/lib/pq" 
)

type Food struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description1 string `json:"description1"`
	Description2 string `json:"description2"`
	Description3 string `json:"description3"`
	Description4 string `json:"description4"`
	ImageURL     string `json:"image_url"`
	Country      string `json:"country"`
}

var db *sql.DB

func init() {
	var err error

	connStr := "user=postgres dbname=apProject password=doiORG2424 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
}

func getFoodByName(w http.ResponseWriter, r *http.Request) {

	foodName := r.URL.Query().Get("name")
	if foodName == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Name parameter is required")
		return
	}

	fmt.Printf("Searching for food: %s\n", foodName)

	query := `SELECT id, food_name, description1, description2, description3, description4, image_url, country 
			  FROM foods WHERE food_name = $1`
	row := db.QueryRow(query, foodName)

	var food Food
	err := row.Scan(&food.ID, &food.Name, &food.Description1, &food.Description2, &food.Description3, &food.Description4, &food.ImageURL, &food.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Couldn't find food: %s\n", foodName) // Сообщение при отсутствии записи
			sendErrorResponse(w, http.StatusNotFound, "Food not found")
		} else {
			fmt.Printf("Database error while searching for food: %s\n", foodName) // Сообщение при ошибке базы данных
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error querying the database: %v", err))
		}
		return
	}

	// Сообщение при успешном нахождении блюда
	fmt.Printf("Found food: %s (ID: %d)\n", food.Name, food.ID)

	w.Header().Set("Content-Type", "application/json")

	// Кодируем структуру food в JSON и отправляем в ответ
	if err := json.NewEncoder(w).Encode(food); err != nil {
		fmt.Printf("Error encoding JSON for food: %s\n", foodName) // Сообщение при ошибке кодирования
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error encoding JSON: %v", err))
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errorResponse)
}

func main() {

	http.HandleFunc("/food", getFoodByName)

	// Cors config
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"}) // Разрешаем все домены

	// Cors HTTP
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", handlers.CORS(origins, methods, headers)(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}
