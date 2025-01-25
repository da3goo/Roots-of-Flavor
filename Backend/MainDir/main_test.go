package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
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
	// Добавляем данные в базу данных через SQL (уже вставил вручную)
	expectedFood := Food{
		ID:           60, // Предполагается, что это правильный ID для вставленных данных
		Name:         "testname",
		Description1: "testdescription1",
		Description2: "testdescription2",
		Description3: "description3",
		Description4: "description4",
		ImageURL:     "testurl",
		Country:      "testcountry",
	}

	// Выполнение запроса
	response, err := http.Get("http://localhost:8080/food?name=testname")
	if err != nil {
		t.Fatalf("Error while calling API: %v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	// Для диагностики
	t.Logf("Response body: %s", body) // Добавлено для вывода тела ответа

	// Разбираем ответ в структуру Food
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
