package main

import (
	"encoding/json"
	"testing"
	"time"
)

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

	// Сериализация в JSON
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to serialize User: %v", err)
	}

	// Проверка десериализации
	var deserializedUser User
	err = json.Unmarshal(data, &deserializedUser)
	if err != nil {
		t.Fatalf("Failed to deserialize User: %v", err)
	}

	// Проверка, что значения совпадают
	if user.ID != deserializedUser.ID || user.Fullname != deserializedUser.Fullname {
		t.Errorf("Expected %v, but got %v", user, deserializedUser)
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

	// Сериализация в JSON
	data, err := json.Marshal(food)
	if err != nil {
		t.Fatalf("Failed to serialize Food: %v", err)
	}

	// Проверка десериализации
	var deserializedFood Food
	err = json.Unmarshal(data, &deserializedFood)
	if err != nil {
		t.Fatalf("Failed to deserialize Food: %v", err)
	}

	// Проверка, что значения совпадают
	if food.Name != deserializedFood.Name || food.Country != deserializedFood.Country {
		t.Errorf("Expected %v, but got %v", food, deserializedFood)
	}
}

func TestRegistrationDataJSONSerialization(t *testing.T) {
	regData := RegistrationData{
		Fullname: "Jane Smith",
		Email:    "janesmith@example.com",
		Password: "securepassword123",
	}

	// Сериализация в JSON
	data, err := json.Marshal(regData)
	if err != nil {
		t.Fatalf("Failed to serialize RegistrationData: %v", err)
	}

	// Проверка десериализации
	var deserializedRegData RegistrationData
	err = json.Unmarshal(data, &deserializedRegData)
	if err != nil {
		t.Fatalf("Failed to deserialize RegistrationData: %v", err)
	}

	// Проверка, что значения совпадают
	if regData.Email != deserializedRegData.Email || regData.Fullname != deserializedRegData.Fullname {
		t.Errorf("Expected %v, but got %v", regData, deserializedRegData)
	}
}
