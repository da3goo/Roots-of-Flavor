package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("key1")) // Секретный ключ для сессий
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type RegisterRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID        int       `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
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

func init() {
	var err error
	connStr := "postgres://postgres.omqkkeruydkttwwkdnib:50TADocqYFe4CFTx@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?sslmode=require&supa=base-pooler.x"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	store = sessions.NewCookieStore([]byte("key1"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login request received.")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to log in user: %s", credentials.Email)

	var user User
	query := "SELECT id, fullname, email FROM users WHERE email = $1 AND password = $2"
	err = db.QueryRow(query, credentials.Email, credentials.Password).Scan(&user.ID, &user.Fullname, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Invalid credentials for user: %s", credentials.Email)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email or password"})
			return
		} else {
			log.Printf("Database error during login: %v", err)
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = user.ID
	session.Values["fullname"] = user.Fullname
	session.Values["email"] = user.Email
	session.Save(r, w)

	log.Printf("User %s logged in successfully", user.Fullname)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout request received.")

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err) // Логирование ошибки при получении сессии
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = nil
	session.Values["fullname"] = nil
	session.Values["email"] = nil
	log.Println("Session values cleared.")

	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err) // Логирование ошибки при сохранении сессии
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("User logged out successfully.")

	// Возвращаем ответ о успешном выходе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func checkSessionHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		log.Println("User not logged in or session expired")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	query := "SELECT id, fullname, email, created_at FROM users WHERE id = $1"
	err = db.QueryRow(query, userID).Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Registration request received.")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received registration request for email: %s", credentials.Email)

	var existingUserID int
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", credentials.Email).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking for existing user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingUserID != 0 {

		log.Printf("User already exists with email: %s", credentials.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"error": "User already exists with this email"})
		return
	}

	log.Println("Inserting new user into database...")
	var newUserID int
	query := "INSERT INTO users (fullname, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = db.QueryRow(query, credentials.Fullname, credentials.Email, credentials.Password).Scan(&newUserID)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("New user registered with ID: %d", newUserID)

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = newUserID
	session.Values["fullname"] = credentials.Fullname
	session.Values["email"] = credentials.Email
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Session created for user: %s", credentials.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func main() {
	http.HandleFunc("/food", getFoodByName)

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/checksession", checkSessionHandler)
	http.HandleFunc("/register", registerHandler)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://127.0.0.1:5500"}) // Указываем конкретный источник
	credentials := handlers.AllowCredentials()

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", handlers.CORS(origins, methods, headers, credentials)(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}

func getFoodByName(w http.ResponseWriter, r *http.Request) {
	foodName := r.URL.Query().Get("name")
	if foodName == "" {
		log.Println("No food name provided.")
		sendErrorResponse(w, http.StatusBadRequest, "Name parameter is required")
		return
	}

	log.Printf("Searching for food: %s", foodName)

	query := `SELECT id, food_name, description1, description2, description3, description4, image_url, country 
			  FROM foods WHERE food_name = $1`
	row := db.QueryRow(query, foodName)

	var food Food
	err := row.Scan(&food.ID, &food.Name, &food.Description1, &food.Description2, &food.Description3, &food.Description4, &food.ImageURL, &food.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Couldn't find food: %s", foodName) // Сообщение при отсутствии записи
			sendErrorResponse(w, http.StatusNotFound, "Food not found")
		} else {
			log.Printf("Database error while searching for food: %s", foodName) // Сообщение при ошибке базы данных
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error querying the database: %v", err))
		}
		return
	}

	log.Printf("Found food: %s (ID: %d)", food.Name, food.ID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(food); err != nil {
		log.Printf("Error encoding JSON for food: %s", foodName) // Сообщение при ошибке кодирования
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error encoding JSON: %v", err))
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errorResponse)
}
