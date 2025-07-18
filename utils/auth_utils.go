package utils

import (
	"github.com/google/uuid"
	"github.com/mahabub618/minipack/internal/repositories"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	sessionCookieName = "session_id"
	sessionLength     = 30 * 24 * time.Hour
)

func GetOrCreateAnonymousUserID(w http.ResponseWriter, r *http.Request, repo *repositories.AnonymousUserRepository) (string, error) {
	// Try to get session ID from cookie
	cookie, err := r.Cookie(sessionCookieName)
	var sessionID string
	if cookie != nil {
		value := cookie.Value
		log.Println("#255 value: ", value)
		if strings.HasPrefix(value, "session_id=") {
			sessionID = strings.TrimPrefix(value, "session_id=")
		} else {
			sessionID = value // fallback in case format changes
		}
	}
	log.Println("#255 request: ", r)
	log.Println("#255 cookie: ", cookie)
	log.Println("#255 error: ", err)

	if err == nil {
		sessionID = cookie.Value
		log.Println("#255 Found existing session ID", sessionID)
	} else {
		// Create new session ID if none exists
		sessionID = uuid.New().String()
		log.Println("#255 Creating new session ID", sessionID)
		// Set the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionID,
			Path:     "/",
			Expires:  time.Now().Add(sessionLength),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		})
	}

	// Check if we already have an anonymous user for this session
	user, err := repo.GetBySessionID(r.Context(), sessionID)
	log.Println("#255 Already have anonymous user", user)
	if err != nil && user != nil {
		return user.ID, err
	}

	if user != nil {
		// Check if expired
		if time.Now().After(user.ExpiresAt) {
			// Create new anonymous user
			newUser, err := repo.Create(r.Context(), sessionID)
			log.Println("#255 Expired user, creating new one", newUser)
			if err != nil {
				return "", err
			}
			return newUser.ID, nil
		}
		log.Println("#255 Final user id", user.ID)
		return user.ID, nil
	}

	// Create new anonymous user
	newUser, err := repo.Create(r.Context(), sessionID)
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
