package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath    = "/path/to/selenium-server.jar"                         // Укажите путь к Selenium Server
	geckoDriverPath = "/path/to/geckodriver"                                 // Укажите путь к WebDriver (например, geckodriver)
	port            = 4444                                                   // Порт для Selenium Server
	baseURL         = "http://127.0.0.1:5500/Frontend/client/main_page.html" // URL вашего приложения
)

// Unit tests

// Structure tests
func TestUserJSONSerialization(t *testing.T) {
	user := User{
		ID:                1,
		Fullname:          "John Doe",
		Email:             "johndoe@example.com",
		CreatedAt:         time.Now(),
		UpdatedFullnameAt: time.Now(),
		Userstatus:        "active",
		UpdatedAt:         time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to serialize User: %v", err)
	}

	var deserializedUser User
	err = json.Unmarshal(data, &deserializedUser)
	if err != nil {
		t.Fatalf("Failed to deserialize User: %v", err)
	}

	if user.ID != deserializedUser.ID || user.Fullname != deserializedUser.Fullname {
		t.Errorf("Expected %v, but got %v", user, deserializedUser)
	}

	if !user.CreatedAt.Equal(deserializedUser.CreatedAt) {
		t.Errorf("Expected CreatedAt: %v, but got %v", user.CreatedAt, deserializedUser.CreatedAt)
	}
	if !user.UpdatedAt.Equal(deserializedUser.UpdatedAt) {
		t.Errorf("Expected UpdatedAt: %v, but got %v", user.UpdatedAt, deserializedUser.UpdatedAt)
	}
}

func TestFoodJSONSerialization(t *testing.T) {
	food := Food{
		ID:           1,
		Name:         "Pizza",
		Description1: "Delicious Italian pizza",
		Description2: "Topped with fresh ingredients",
		Description3: "",
		Description4: "",
		ImageURL:     "http://example.com/pizza.jpg",
		Country:      "Italy",
	}

	data, err := json.Marshal(food)
	if err != nil {
		t.Fatalf("Failed to serialize Food: %v", err)
	}

	var deserializedFood Food
	err = json.Unmarshal(data, &deserializedFood)
	if err != nil {
		t.Fatalf("Failed to deserialize Food: %v", err)
	}

	if food.Name != deserializedFood.Name || food.Country != deserializedFood.Country {
		t.Errorf("Expected %v, but got %v", food, deserializedFood)
	}

	if food.Description1 != deserializedFood.Description1 {
		t.Errorf("Expected Description1: %v, but got %v", food.Description1, deserializedFood.Description1)
	}
	if food.Description2 != deserializedFood.Description2 {
		t.Errorf("Expected Description2: %v, but got %v", food.Description2, deserializedFood.Description2)
	}
	if food.Description3 != deserializedFood.Description3 {
		t.Errorf("Expected Description3: %v, but got %v", food.Description3, deserializedFood.Description3)
	}
	if food.Description4 != deserializedFood.Description4 {
		t.Errorf("Expected Description4: %v, but got %v", food.Description4, deserializedFood.Description4)
	}
}

func TestRegistrationDataJSONSerialization(t *testing.T) {
	regData := RegistrationData{
		Fullname: "Jane Smith",
		Email:    "janesmith@example.com",
		Password: "securepassword123",
	}

	data, err := json.Marshal(regData)
	if err != nil {
		t.Fatalf("Failed to serialize RegistrationData: %v", err)
	}

	var deserializedRegData RegistrationData
	err = json.Unmarshal(data, &deserializedRegData)
	if err != nil {
		t.Fatalf("Failed to deserialize RegistrationData: %v", err)
	}

	if regData.Email != deserializedRegData.Email || regData.Fullname != deserializedRegData.Fullname {
		t.Errorf("Expected %v, but got %v", regData, deserializedRegData)
	}
}

func TestGetFoodByName(t *testing.T) {

	expectedFood := Food{
		ID:           60,
		Name:         "testname",
		Description1: "testdescription1",
		Description2: "testdescription2",
		Description3: "description3",
		Description4: "description4",
		ImageURL:     "testurl",
		Country:      "testcountry",
	}

	response, err := http.Get("http://localhost:8080/food?name=testname")
	if err != nil {
		t.Fatalf("Error while calling API: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	t.Logf("Response body: %s", body)

	var actualFood Food
	err = json.Unmarshal(body, &actualFood)
	if err != nil {
		t.Fatalf("Error unmarshaling response: %v", err)
	}

	// Сравниваем ожидаемые и фактические данные
	if expectedFood != actualFood {
		t.Errorf("Expected %v, but got %v", expectedFood, actualFood)
	}
}

// Selenium testing

func TestLoginWithOTP(t *testing.T) {
	const (
		chromeDriverPath = "./chromedriver-win64/chromedriver.exe"
		baseURL          = "http://127.0.0.1:5500/Frontend/client/main_page.html"
	)

	// ChromeDriver settings
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444)
	if err != nil {
		t.Fatalf("Error starting the ChromeDriver service: %v", err)
	}
	defer service.Stop()

	// WebDriver settings for Chrome
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:62977")
	if err != nil {
		t.Fatalf("Error connecting to the WebDriver: %v", err)
	}
	defer wd.Quit()

	// Main page load
	if err := wd.Get(baseURL); err != nil {
		t.Fatalf("Failed to load page: %v", err)
	}

	// Open modal page
	loginLink, err := wd.FindElement(selenium.ByID, "logInText")
	if err != nil {
		t.Fatalf("Failed to find login link: %v", err)
	}
	loginLink.Click()

	//Modal page waiting
	time.Sleep(1 * time.Second)

	// Email input
	emailInput, err := wd.FindElement(selenium.ByID, "emailInput")
	if err != nil {
		t.Fatalf("Failed to find email input: %v", err)
	}
	emailInput.SendKeys("kantaidaulet@gmail.com")

	// Password Input
	passwordInput, err := wd.FindElement(selenium.ByID, "passwordInput")
	if err != nil {
		t.Fatalf("Failed to find password input: %v", err)
	}
	passwordInput.SendKeys("qwerty1234")

	// Clicking to Login
	loginButton, err := wd.FindElement(selenium.ByCSSSelector, ".modal-button")
	if err != nil {
		t.Fatalf("Failed to find login button: %v", err)
	}
	loginButton.Click()

	time.Sleep(3 * time.Second)

	// DB connection for otp code
	db, err := sql.Open("postgres", "postgres://postgres.omqkkeruydkttwwkdnib:50TADocqYFe4CFTx@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?sslmode=require&supa=base-pooler.x") // Укажите строку подключения к вашей БД
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	var otp string
	query := "SELECT otp FROM users WHERE email = $1 AND otp IS NOT NULL"
	err = db.QueryRow(query, "kantaidaulet@gmail.com").Scan(&otp)
	if err != nil {
		t.Fatalf("Failed to fetch OTP: %v", err)
	}

	// OTP entering
	otpInput, err := wd.FindElement(selenium.ByID, "verificationCodeLogin")
	if err != nil {
		t.Fatalf("Failed to find OTP input: %v", err)
	}
	otpInput.SendKeys(otp)

	// Verify click
	verifyButton, err := wd.FindElement(selenium.ByID, "verifyCodeButtonLogin")
	if err != nil {
		t.Fatalf("Failed to find verify button: %v", err)
	}
	verifyButton.Click()

	time.Sleep(3 * time.Second)

	time.Sleep(3 * time.Second)

	// Checking
	profileLink, err := wd.FindElement(selenium.ByLinkText, "Profile")
	if err != nil {
		t.Fatalf("Profile link not found after reload: %v", err)
	}

	t.Log("Profile link:", profileLink)
	t.Log("Login successful!")
}
