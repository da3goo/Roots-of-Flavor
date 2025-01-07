package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

var (
	db      *sql.DB
	store   = sessions.NewCookieStore([]byte("key1"))
	limiter = rate.NewLimiter(1, 3)
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
	http.HandleFunc("/changePassword", changePassword)
	http.HandleFunc("/changeEmail", changeEmail)
	http.HandleFunc("/send", handleForm)

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
	UpdatedAt         time.Time `json:"updatedAt"`
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

	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.SetOutput(os.Stdout)

	var err error
	connStr := "postgres://postgres.omqkkeruydkttwwkdnib:50TADocqYFe4CFTx@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?sslmode=require&supa=base-pooler.x"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"connStr": connStr,
		}).Fatal("Failed to open database connection")
	}

	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error connecting to the database")
	} else {
		logrus.WithFields(logrus.Fields{
			"status": "success",
		}).Info("Connected to Database")
	}
}
func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func login(w http.ResponseWriter, r *http.Request) {
	// Rate limiting
	if !limiter.Allow() {
		resetTime := time.Now().Add(time.Second * time.Duration(limiter.Reserve().Delay()))

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(resetTime.Unix())))

		http.Error(w, "Rate limit exceeded, try again later", http.StatusTooManyRequests)
		return
	}

	logrus.WithFields(logrus.Fields{
		"method":   r.Method,
		"endpoint": "/login",
	}).Info("Request started")

	if r.Method != http.MethodPost {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"status": "fail",
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"status": "fail",
		}).Error("Failed to decode JSON")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid JSON format",
			"status":  "fail",
		})
		return
	}

	var user User
	var storedPassword string // Temporary variable for the password
	// Query to get the password hash
	query := `
        SELECT id, fullname, email, password, updated_fullname_at, userstatus 
        FROM users 
        WHERE email = $1`
	err = db.QueryRow(query, credentials.Email).Scan(
		&user.ID, &user.Fullname, &user.Email, &storedPassword, &user.UpdatedFullnameAt, &user.Userstatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithFields(logrus.Fields{
				"email":  credentials.Email,
				"status": "fail",
			}).Warn("Invalid email or password")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid email or password",
				"status":  "fail",
			})
			return
		}
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"status": "fail",
		}).Error("Database error during login")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"email":  credentials.Email,
			"status": "fail",
		}).Warn("Invalid email or password")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid email or password",
			"status":  "fail",
		})
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID": user.ID,
		"email":  user.Email,
		"status": "success",
	}).Info("User logged in successfully")

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

func logout(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"method":   r.Method,
		"endpoint": "/logout",
	}).Info("Logout request received")

	session, err := store.Get(r, "user-session")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error retrieving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.Info("Session values cleared.")
	session.Values["userID"] = nil
	session.Values["fullname"] = nil
	session.Values["email"] = nil

	err = session.Save(r, w)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error saving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.Info("User logged out successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func checkSession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")

	logrus.WithFields(logrus.Fields{
		"session": session.Values,
	}).Info("Session contents")

	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		logrus.WithFields(logrus.Fields{
			"session": session.Values,
		}).Warn("The session is invalid or the userID is missing")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID": userID,
	}).Info("The user ID was found in the session")

	var user User
	query := `
        SELECT id, fullname, email, created_at, updated_fullname_at, updated_at, userstatus 
        FROM users 
        WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt, &user.UpdatedFullnameAt, &user.UpdatedAt, &user.Userstatus)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error when requesting a user from the database")
		if err == sql.ErrNoRows {
			logrus.WithFields(logrus.Fields{
				"userID": userID,
			}).Warn("No rows found for the user ID")
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID":   user.ID,
		"fullname": user.Fullname,
	}).Info("User is found")

	// Добавляем updatedAt в ответ
	response := map[string]interface{}{
		"id":                user.ID,
		"fullname":          user.Fullname,
		"email":             user.Email,
		"createdAt":         user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedFullnameAt": user.UpdatedFullnameAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt":         user.UpdatedAt.Format("2006-01-02T15:04:05Z"), // добавляем updatedAt
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

	// Проверка существующего пользователя с таким же email
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

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Inserting new user into database...")
	var newUserID int
	query := "INSERT INTO users (fullname, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = db.QueryRow(query, credentials.Fullname, credentials.Email, string(hashedPassword)).Scan(&newUserID)
	if err != nil {
		log.Printf("Error inserting new user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("New user registered with ID: %d", newUserID)

	// Создаем сессию
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
	logrus.Info("Delete profile request received.")

	if r.Method != http.MethodDelete {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error retrieving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		logrus.Warn("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(query, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Error deleting user")
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error deleting session")
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	logrus.WithFields(logrus.Fields{
		"userID": userID,
	}).Info("User deleted successfully")
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

	logrus.WithFields(logrus.Fields{
		"emailFilter": emailFilter,
		"sortBy":      sortBy,
		"page":        page,
		"pageSize":    pageSize,
	}).Info("Received request to get users with filters")

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
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"query": query,
		}).Error("Error executing query")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Fullname, &user.Email, &user.CreatedAt, &user.UpdatedFullnameAt, &user.Userstatus); err != nil {
			http.Error(w, fmt.Sprintf("Error when reading data: %v", err), http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error when reading line data")
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error in processing the results: %v", err), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error processing rows")
		return
	}

	if len(users) == 0 {
		http.Error(w, "No users found matching the filters", http.StatusNotFound)
		logrus.Info("No users found matching the filters")
		return
	}

	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email LIKE $1", "%"+emailFilter+"%").Scan(&totalCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting total count: %v", err), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error getting total count")
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
	logrus.WithFields(logrus.Fields{
		"userCount": len(users),
	}).Info("Sending response with user data")

	json.NewEncoder(w).Encode(response)
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	//Request limitting
	if !limiter.Allow() {
		resetTime := time.Now().Add(time.Second * time.Duration(limiter.Reserve().Delay()))

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(resetTime.Unix())))

		logrus.WithFields(logrus.Fields{
			"resetTime": resetTime,
			"remaining": limiter.Burst(),
			"limit":     limiter.Limit(),
		}).Warn("Rate limit exceeded")

		http.Error(w, "Rate limit exceeded, try again later", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		OldPassword       string `json:"oldPassword"`
		NewPassword       string `json:"newPassword"`
		NewPasswordRetype string `json:"newPasswordRetype"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid JSON format")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error retrieving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"sessionValues": session.Values,
	}).Debug("Session values retrieved")

	userID, ok := session.Values["userID"].(int)
	if session.Values["userID"] == nil {
		logrus.Warn("Session is missing userID")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	if !ok {
		logrus.Warn("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	var currentPassword string
	query := `SELECT password FROM users WHERE id = $1`
	err = db.QueryRow(query, userID).Scan(&currentPassword)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Error retrieving user password")
		http.Error(w, "Failed to retrieve user password", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(currentPassword), []byte(requestData.OldPassword))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Warn("Incorrect old password")
		http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
		return
	}

	if requestData.NewPassword != requestData.NewPasswordRetype {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
		}).Warn("New passwords do not match")
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Error hashing new password")
		http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
		return
	}

	updateQuery := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`
	_, err = db.Exec(updateQuery, string(hashedPassword), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":   userID,
			"error":    err.Error(),
			"password": requestData.NewPassword,
		}).Error("Error updating user password")
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	session.Values["updatedAt"] = time.Now().Format("2006-01-02 15:04:05")
	err = session.Save(r, w)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Error saving session after password update")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID": userID,
	}).Info("Password updated successfully")

	response := map[string]string{
		"message":    "Password updated successfully",
		"updated_at": session.Values["updatedAt"].(string),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func changeEmail(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		resetTime := time.Now().Add(time.Second * time.Duration(limiter.Reserve().Delay()))

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Limit()))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(resetTime.Unix())))

		logrus.WithFields(logrus.Fields{
			"resetTime": resetTime,
			"remaining": limiter.Burst(),
			"limit":     limiter.Limit(),
		}).Warn("Rate limit exceeded")

		http.Error(w, "Rate limit exceeded, try again later", http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		NewEmail string `json:"newEmail"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid JSON format")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error retrieving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		logrus.Warn("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	var existingUserID int
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", requestData.NewEmail).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithFields(logrus.Fields{
			"newEmail": requestData.NewEmail,
			"error":    err.Error(),
		}).Error("Database error when checking email existence")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingUserID != 0 {
		logrus.WithFields(logrus.Fields{
			"newEmail": requestData.NewEmail,
		}).Warn("Email already exists")
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	_, err = db.Exec("UPDATE users SET email = $1 WHERE id = $2", requestData.NewEmail, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID":   userID,
			"newEmail": requestData.NewEmail,
			"error":    err.Error(),
		}).Error("Failed to update email")
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}

	session.Values["email"] = requestData.NewEmail
	err = session.Save(r, w)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err.Error(),
		}).Error("Failed to save session after email update")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logrus.WithFields(logrus.Fields{
		"userID":   userID,
		"newEmail": requestData.NewEmail,
	}).Info("Email updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Email updated successfully"})
}

func sendEmail(from, password, to, subject, message, filename, fileContent string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	logrus.WithFields(logrus.Fields{
		"from":    from,
		"to":      to,
		"subject": subject,
	}).Info("Sending email started")

	mime := "MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=\"boundary1\"\n\n"
	body := "--boundary1\n"
	body += "Content-Type: text/plain; charset=\"utf-8\"\n\n"
	body += message + "\n\n"

	if filename != "" && fileContent != "" {
		body += "--boundary1\n"
		body += "Content-Type: application/octet-stream; name=\"" + filename + "\"\n"
		body += "Content-Disposition: attachment; filename=\"" + filename + "\"\n"
		body += "Content-Transfer-Encoding: base64\n\n"
		body += fileContent + "\n"
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		mime + body + "--boundary1--"

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"from":    from,
			"to":      to,
			"subject": subject,
			"error":   err.Error(),
		}).Error("Failed to send email")
		return err
	}

	logrus.WithFields(logrus.Fields{
		"from":    from,
		"to":      to,
		"subject": subject,
	}).Info("Email sent successfully")

	return nil
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		logrus.WithFields(logrus.Fields{
			"email":   email,
			"subject": subject,
			"message": message,
		}).Info("Processing form submission")

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error receiving file")
			http.Error(w, "Error receiving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		logrus.WithFields(logrus.Fields{
			"filename": fileHeader.Filename,
			"size":     fileHeader.Size,
		}).Info("File uploaded")

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error reading file")
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		fileContent := base64.StdEncoding.EncodeToString(fileBytes)

		_, filename := filepath.Split(fileHeader.Filename)

		err = sendEmail(email, password, "kantaydaulet777@gmail.com", subject, message, filename, fileContent)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"from":    email,
				"to":      "kantaydaulet777@gmail.com",
				"subject": subject,
				"error":   err.Error(),
			}).Error("Error sending email")
			http.Error(w, "Error sending email: "+err.Error(), http.StatusInternalServerError)
			return
		}

		logrus.WithFields(logrus.Fields{
			"from":    email,
			"to":      "kantaydaulet777@gmail.com",
			"subject": subject,
		}).Info("Email sent successfully")

		fmt.Fprintf(w, "The message was sent successfully!")
	} else {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
		}).Warn("Unsupported method")
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}
