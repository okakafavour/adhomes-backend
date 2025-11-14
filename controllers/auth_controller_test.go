package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"adhomes-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------
// TestMain - runs before any test
// -------------------------------
func TestMain(m *testing.M) {
	// Load .env
	if err := godotenv.Load("../.env"); err != nil {
		panic("❌ Error loading .env file")
	}

	// Connect to test database
	config.ConnectTestDB()

	// Initialize controllers
	InitUserController()
	InitProductController()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

// -------------------------------
// Helpers
// -------------------------------

// Clean users collection safely
func cleanUsersCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.DB.Collection("users").Drop(ctx); err != nil {
		panic(err)
	}
}

// Setup router
func setUpUserRouter() *gin.Engine {
	r := gin.New()
	r.POST("/signup", SignUp)
	r.POST("/login", Login)
	return r
}

// -------------------------------
// Tests
// -------------------------------

func TestToRegisterUser(t *testing.T) {
	cleanUsersCollection()
	router := setUpUserRouter()

	body := []byte(`{
		"name": "Okaka Favour",
		"email": "okaka@example.com",
		"password": "12345"
	}`)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "User created successfully", resp["message"])
}

func TestToRegisterWithoutEmail(t *testing.T) {
	cleanUsersCollection()
	router := setUpUserRouter()

	body := []byte(`{"name": "Sam", "password": "0000"}`)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Email is required", resp["error"])
}

func TestToLoginUser(t *testing.T) {
	cleanUsersCollection()
	router := setUpUserRouter()

	// First, register user
	signupBody := []byte(`{
		"name": "Okaka Favour",
		"email": "okaka@example.com",
		"password": "12345"
	}`)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Now login
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Login successful", resp["message"])
	assert.NotEmpty(t, resp["token"])
}
