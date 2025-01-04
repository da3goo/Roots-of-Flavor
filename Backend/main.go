package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	db    *sql.DB
	store = sessions.NewCookieStore([]byte("key1"))
)

func main() {
	http.HandleFunc("/food", getFoodByName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/checksession", checkSession)
	http.HandleFunc("/register", register)
	http.HandleFunc("/updateName", updateName)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/getUsers", getUsers)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://127.0.0.1:5500"})
	credentials := handlers.AllowCredentials()

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", handlers.CORS(origins, methods, headers, credentials)(http.DefaultServeMux))
	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}

type User struct {
	ID                int       `json:"id"`
	Fullname          string    `json:"fullname"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedFullnameAt time.Time `json:"updatedFullnameAt"`
	Userstatus        string    `json:"userstatus"`
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

// Backend functions
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
	} else {
		log.Println("Connected to Database")
	}

}
func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON format",
			"status":  "fail",
		})
		return
	}

	var user User
	query := `
		SELECT id, fullname, email, updated_fullname_at, userstatus 
		FROM users 
		WHERE email = $1 AND password = $2`
	err = db.QueryRow(query, credentials.Email, credentials.Password).Scan(
		&user.ID, &user.Fullname, &user.Email, &user.UpdatedFullnameAt, &user.Userstatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid email or password",
				"status":  "fail",
			})
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "user-session")
	session.Values["userID"] = user.ID
	session.Values["fullname"] = user.Fullname
	session.Values["email"] = user.Email
	session.Values["updatedFullnameAt"] = user.UpdatedFullnameAt.Format("2006-01-02 15:04:05")
	session.Values["userstatus"] = user.Userstatus
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 32
func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout request received.")

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["userID"] = nil
	session.Values["fullname"] = nil
	session.Values["email"] = nil
	log.Println("Session values cleared.")

	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("User logged out successfully.")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func checkSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")

	// Логируем всю сессию
	log.Printf("Session contents: %v", session.Values)

	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		log.Println("The session is invalid or the userID is missing")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("The user ID was found in the session: %d", userID)

	var user User
	query := "SELECT id, fullname, email, created_at, updated_fullname_at, userstatus FROM users WHERE id = $1"
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt, &user.UpdatedFullnameAt, &user.Userstatus)
	if err != nil {
		// Логируем ошибку при запросе пользователя
		log.Printf("Error when requesting a user from the database: %v", err)
		if err == sql.ErrNoRows {
			log.Println("No rows found for the user ID")
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("User is found %v", user)

	response := map[string]interface{}{
		"id":                user.ID,
		"fullname":          user.Fullname,
		"email":             user.Email,
		"createdAt":         user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedFullnameAt": user.UpdatedFullnameAt.Format("2006-01-02T15:04:05Z"),
		"userStatus":        user.Userstatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func register(w http.ResponseWriter, r *http.Request) {
	log.Println("Registration request received.")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var rawData map[string]interface{}
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	expectedKeys := map[string]string{
		"fullname": "string",
		"email":    "string",
		"password": "string",
	}

	for key, expectedType := range expectedKeys {
		value, exists := rawData[key]
		if !exists {
			log.Printf("Missing key: %s", key)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Incorrect key name"})
			return
		}

		switch expectedType {
		case "string":
			if _, ok := value.(string); !ok {
				log.Printf("Invalid type for key: %s, expected: %s", key, expectedType)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Invalid type for key: " + key})
				return
			}
		}
	}

	var credentials struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Printf("Error decoding JSON into struct: %v", err)
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

func updateName(w http.ResponseWriter, r *http.Request) {
	log.Println("Profile update request received.")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Fullname string `json:"fullname"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		log.Println("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	query := `UPDATE users 
              SET fullname = $1, updated_fullname_at = NOW() 
              WHERE id = $2 
              RETURNING updated_fullname_at`
	var updatedAt time.Time
	err = db.QueryRow(query, requestData.Fullname, userID).Scan(&updatedAt)
	if err != nil {
		log.Printf("Error updating user profile: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":             "Profile updated successfully",
		"updated_fullname_at": updatedAt.Format("2006-01-02 15:04:05"),
	}
	json.NewEncoder(w).Encode(response)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete profile request received.")

	if r.Method != http.MethodDelete {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		log.Println("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(query, userID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
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
			log.Printf("Couldn't find food: %s", foodName)
			sendErrorResponse(w, http.StatusNotFound, "Food not found")
		} else {
			log.Printf("Database error while searching for food: %s", foodName)
			sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error querying the database: %v", err))
		}
		return
	}

	log.Printf("Found food: %s (ID: %d)", food.Name, food.ID)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(food); err != nil {
		log.Printf("Error encoding JSON for food: %s", foodName)
		sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error encoding JSON: %v", err))
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{"error": message}
	json.NewEncoder(w).Encode(errorResponse)
}
func getUsers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	emailFilter := r.URL.Query().Get("email")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	log.Printf("Received request to get users with filters - Email filter: %s, Sort by: %s", emailFilter, sortBy)

	pageInt := 1
	pageSizeInt := 9
	if page != "" {
		fmt.Sscanf(page, "%d", &pageInt)
	}
	if pageSize != "" {
		fmt.Sscanf(pageSize, "%d", &pageSizeInt)
	}

	query := "SELECT id, fullname, email, created_at, updated_fullname_at, userstatus FROM users WHERE email LIKE $1"
	var args []interface{}
	args = append(args, "%"+emailFilter+"%")

	switch sortBy {
	case "nameAsc":
		query += " ORDER BY fullname ASC"
	case "nameDesc":
		query += " ORDER BY fullname DESC"
	case "createdAt":
		query += " ORDER BY created_at"
	case "id":
		query += " ORDER BY id"
	default:
		query += " ORDER BY created_at DESC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSizeInt, (pageInt-1)*pageSizeInt)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt, &user.UpdatedFullnameAt, &user.Userstatus); err != nil {
			http.Error(w, fmt.Sprintf("Error when reading data: %v", err), http.StatusInternalServerError)
			log.Printf("Error when reading line data: %v", err)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error in processing the results: %v", err), http.StatusInternalServerError)
		log.Printf("Error processing rows: %v", err)
		return
	}

	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email LIKE $1", "%"+emailFilter+"%").Scan(&totalCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting total count: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting total count: %v", err)
		return
	}

	totalPages := (totalCount + pageSizeInt - 1) / pageSizeInt

	response := map[string]interface{}{
		"users":       users,
		"totalCount":  totalCount,
		"totalPages":  totalPages,
		"currentPage": pageInt,
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("Sending response with %d users", len(users))
	json.NewEncoder(w).Encode(response)
}
