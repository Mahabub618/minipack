package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mahabub618/minipack/config"
	"github.com/mahabub618/minipack/internal/models"
	"github.com/mahabub618/minipack/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

var ErrEmailExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

// RegisterUser handles user registration
func (s *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	exists, err := s.userRepo.EmailExists(ctx, user.Email)
	if err != nil {
		return errors.New("database error")
	}
	if exists {
		return ErrEmailExists
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword
	user.Role = "user"
	user.IsVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, string, error) {
	// Retrieve user by email
	dbUser, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	if !CheckPasswordHash(password, dbUser.Password) {
		return "", "", ErrInvalidCredentials
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", "", errors.New("configuration error")
	}

	// Generate access token
	accessToken, err := GenerateAccessToken(dbUser.ID, dbUser.Role, cfg)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := GenerateRefreshToken(dbUser.ID, cfg)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
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

func GenerateRefreshToken(userID int, cfg *config.Config) (string, error) {
	expirationTime := time.Now().Add(cfg.RefreshTokenExpiry)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = expirationTime.Unix()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	return tokenString, err
}
