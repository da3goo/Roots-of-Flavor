package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
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
	http.HandleFunc("/deleteUserAdmin", deleteUserFromAdminPage)
	http.HandleFunc("/verifyCode", verifyCode)
	http.HandleFunc("/verifyOTP", verifyOTP)

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
type RegistrationData struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	desktopPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Cannot get desktop directory: %v", err)
	}
	logFilePath := filepath.Join(desktopPath, "Desktop", "logs.txt")

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Couldn't open log file: %v", err)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(os.Stdout, logFile))

	connStr := "postgres://postgres.omqkkeruydkttwwkdnib:50TADocqYFe4CFTx@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?sslmode=require&supa=base-pooler.x"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"connStr": connStr,
		}).Fatal("Error connecting to DB")
	}

	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error connecting to DB")
	} else {
		logrus.WithFields(logrus.Fields{
			"status": "success",
		}).Info("Successfully connected to DB")
	}
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func login(w http.ResponseWriter, r *http.Request) {
	// Rate limiting

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
	var storedPassword string
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

	// Генерация OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000)) // 6-значный код
	expiry := time.Now().Add(5 * time.Minute)      // Срок действия OTP: 5 минут

	// Сохраняем OTP и время его действия в базе данных
	_, err = db.Exec("UPDATE users SET otp = $1, otp_expiry = $2 WHERE id = $3", otp, expiry, user.ID)
	if err != nil {
		logrus.Error("Failed to save OTP:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправка OTP на email
	go sendEmail("kantaidaulet@gmail.com", "vxaf gbyk lqqy zhyb", user.Email, "Your OTP Code", fmt.Sprintf("Your OTP is: %s", otp), "", "")

	//Cookies
	session, _ := store.Get(r, "user-session")
	session.Values["userID"] = user.ID
	session.Values["fullname"] = user.Fullname
	session.Values["email"] = user.Email
	session.Values["updatedFullnameAt"] = user.UpdatedFullnameAt.Format("2006-01-02 15:04:05")
	session.Values["userstatus"] = user.Userstatus
	session.Save(r, w)

	// Отправляем успешный ответ с данными пользователя (не с OTP)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	logrus.WithFields(logrus.Fields{
		"userID": user.ID,
		"email":  user.Email,
		"status": "success",
	}).Info("User logged in successfully")

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

	// Логируем содержимое сессии для диагностики
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

	response := map[string]interface{}{
		"id":                user.ID,
		"fullname":          user.Fullname,
		"email":             user.Email,
		"createdAt":         user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updatedFullnameAt": user.UpdatedFullnameAt.Format("2006-01-02T15:04:05Z"),
		"updatedAt":         user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		"userStatus":        user.Userstatus,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to encode response to JSON")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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

	var credentials RegistrationData
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		log.Printf("Error decoding JSON into struct: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Received registration request for email: %s", credentials.Email)

	if credentials.Fullname == "" || credentials.Email == "" || credentials.Password == "" {
		log.Println("Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var existingUserID int
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", credentials.Email).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking for existing user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingUserID != 0 {
		log.Printf("User already exists with email: %s", credentials.Email)
		http.Error(w, "User already exists with this email", http.StatusConflict)
		return
	}

	verificationCode := fmt.Sprintf("%04d", rand.Intn(10000))

	var newTempUserID int
	query := `
		INSERT INTO temp_users (fullname, email, password, verification_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) 
		RETURNING id`
	err = db.QueryRow(query, credentials.Fullname, credentials.Email, credentials.Password, verificationCode).Scan(&newTempUserID)
	if err != nil {
		log.Printf("Error inserting temp user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Printf("Temporary user created with ID: %d", newTempUserID)

	go func() {
		subject := "Your Verification Code"
		message := fmt.Sprintf("Here is your verification code: %s", verificationCode)
		from := "kantaidaulet@gmail.com"
		password := "vxaf gbyk lqqy zhyb"

		err := sendEmail(from, password, credentials.Email, subject, message, "", "")
		if err != nil {
			log.Printf("Failed to send verification email to %s: %v", credentials.Email, err)
		} else {
			log.Printf("Verification code sent to %s: %s", credentials.Email, verificationCode)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification code sent to email"})
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
	// Rate limiting
	if !limiter.Allow() {
		resetTime := time.Now().Add(time.Second * time.Duration(limiter.Reserve().Delay()))

		// Удаляем использование limiter.Limit() и вместо этого можем не устанавливать заголовки или установить фиксированные значения
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Burst()))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(resetTime.Unix())))

		logrus.WithFields(logrus.Fields{
			"resetTime": resetTime,
			"remaining": limiter.Burst(),
		}).Warn("Rate limit exceeded")

		http.Error(w, "Rate limit exceeded, try again later", http.StatusTooManyRequests)
		return
	}

	// Обработка метода запроса
	if r.Method != http.MethodPost {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
		}).Warn("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Декодирование JSON тела запроса
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

	// Получение сессии
	session, err := store.Get(r, "user-session")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error retrieving session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Получение userID из сессии
	userID, ok := session.Values["userID"].(int)
	if !ok || userID == 0 {
		logrus.Warn("User not logged in")
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Проверка существования email в базе данных
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

	// Если email уже существует
	if existingUserID != 0 {
		logrus.WithFields(logrus.Fields{
			"newEmail": requestData.NewEmail,
		}).Warn("Email already exists")
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Обновление email в базе данных
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

	// Обновление email в сессии
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

	// Ответ клиенту
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

		var filename string
		var fileContent string

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			if err == http.ErrMissingFile {
				logrus.Info("No file uploaded, proceeding without attachment")
			} else {
				logrus.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Error receiving file")
				http.Error(w, "Error receiving file", http.StatusInternalServerError)
				return
			}
		} else {
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

			fileContent = base64.StdEncoding.EncodeToString(fileBytes)
			_, filename = filepath.Split(fileHeader.Filename)
		}

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

func deleteUserFromAdminPage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User with ID %s deleted successfully", id)
}

func verifyCode(w http.ResponseWriter, r *http.Request) {
	log.Println("Verification request received.")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method:", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Email            string `json:"email"`
		VerificationCode string `json:"code"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var tempUserID int
	var storedCode, plainPassword string
	err = db.QueryRow("SELECT id, verification_code, password FROM temp_users WHERE email = $1", requestData.Email).Scan(&tempUserID, &storedCode, &plainPassword)
	if err != nil {
		log.Printf("Error retrieving temp user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if storedCode != requestData.VerificationCode {
		log.Println("Invalid verification code")
		http.Error(w, "Invalid verification code", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var newUserID int
	query := `
		INSERT INTO users (fullname, email, password, userstatus, created_at, updated_at)
		SELECT fullname, email, $1, 'active', NOW(), NOW()
		FROM temp_users WHERE id = $2
		RETURNING id`
	err = db.QueryRow(query, hashedPassword, tempUserID).Scan(&newUserID)
	if err != nil {
		log.Printf("Error inserting user into main table: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM temp_users WHERE id = $1", tempUserID)
	if err != nil {
		log.Printf("Error deleting temp user: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	log.Println("User successfully verified and moved to main table")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified and user registered"})
}

func verifyOTP(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	var storedOTP string
	var otpExpiry time.Time
	var userID int

	err := db.QueryRow("SELECT id, otp, otp_expiry FROM users WHERE email = $1", input.Email).Scan(&userID, &storedOTP, &otpExpiry)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if storedOTP != input.OTP || time.Now().After(otpExpiry) {
		http.Error(w, "Invalid or expired OTP", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("UPDATE users SET otp = NULL, otp_expiry = NULL WHERE id = $1", userID)
	if err != nil {
		http.Error(w, "Failed to clear OTP", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP verified"})
}
