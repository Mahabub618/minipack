package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/services"
	"net/http"
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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

var validate = validator.New()

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	// Create a new user model
	newUser := models.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		AddressStreet:  sql.NullString{String: req.AddressStreet, Valid: req.AddressStreet != ""},
		AddressCity:    sql.NullString{String: req.AddressCity, Valid: req.AddressCity != ""},
		AddressCountry: sql.NullString{String: req.AddressCountry, Valid: req.AddressCountry != ""},
		Role:           "user",
	}

	// Register the user using the service
	if err := h.userService.RegisterUser(r.Context(), &newUser); err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	//// Get user from DB
	//var user models.User
	//err := database.DB.Get(&user, "SELECT id, email, password, role FROM users WHERE email = $1", req.Email)
	//if err == sql.ErrNoRows {
	//	http.Error(w, "User not found", http.StatusNotFound)
	//	return
	//} else if err != nil {
	//	http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Validate the password
	//if !CheckPasswordHash(req.Password, user.Password) {
	//	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	//	return
	//}
	//
	//// Load config
	//cfg, err := config.LoadConfig()
	//if err != nil {
	//	http.Error(w, "Configuration error", http.StatusInternalServerError)
	//	return
	//}
	//
	//// Generate access token
	//accessToken, err := GenerateAccessToken(user.ID, user.Role, cfg)
	//if err != nil {
	//	http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
	//	return
	//}
	//
	//// Generate refresh token
	//refreshToken, err := GenerateRefreshToken(user.ID, cfg)
	//if err != nil {
	//	http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
	//	return
	//}

	//// Respond with tokens
	//response := LoginResponse{
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(response)
}
