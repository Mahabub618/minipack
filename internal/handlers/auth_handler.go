package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/mahabub618/minipack/internal/database"
	"github.com/mahabub618/minipack/internal/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/config"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type SignUpRequest struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=8"`
	Phone          string `json:"phone"`
	AddressStreet  string `json:"address_street"`
	AddressCity    string `json:"address_city"`
	AddressCountry string `json:"address_country"`
}

type SignupResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone,omitempty"`
	AddressStreet  string `json:"address_street,omitempty"`
	AddressCity    string `json:"address_city,omitempty"`
	AddressCountry string `json:"address_country,omitempty"`
	Role           string `json:"role"`
}

var validate = validator.New()

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	// Decode body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request body
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if email already exists
	var exists bool
	err := database.DB.Get(&exists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user
	newUser := models.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       hashedPassword,
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		AddressStreet:  sql.NullString{String: req.AddressStreet, Valid: req.AddressStreet != ""},
		AddressCity:    sql.NullString{String: req.AddressCity, Valid: req.AddressCity != ""},
		AddressCountry: sql.NullString{String: req.AddressCountry, Valid: req.AddressCountry != ""},
		Role:           "user",
	}

	// Insert user

	query := `
		INSERT INTO users (name, email, password, phone, address_street, address_city, address_country, role)
		VALUES(:name, :email, :password, :phone, :address_street, :address_city, :address_country, :role)
		RETURNING id`

	stmt, err := database.DB.PrepareNamed(query)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.Get(&newUser.ID, newUser)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := SignupResponse{
		ID:             newUser.ID,
		Name:           newUser.Name,
		Email:          newUser.Email,
		Role:           newUser.Role,
		Phone:          newUser.Phone.String,
		AddressStreet:  newUser.AddressStreet.String,
		AddressCity:    newUser.AddressCity.String,
		AddressCountry: newUser.AddressCountry.String,
	}

	// Respond with created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccessToken(userID int, role string, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(cfg.AccessTokenExpiry)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func GenerateRefreshToken(userID int, cfg *config.Config) (string, time.Time, error) {
	expirationTime := time.Now().Add(cfg.RefreshTokenExpiry)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = expirationTime.Unix()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	return tokenString, expirationTime, err
}
