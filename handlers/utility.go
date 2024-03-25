package handlers

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func generateToken(userID int64, userEmail string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"id":    userID,
		"email": userEmail,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":   time.Now().Unix(),                     // Issued at time
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getUserIDFromToken adalah fungsi untuk mendapatkan user ID dari token JWT
func getUserIDFromToken(r *http.Request) (uint, error) {
	// Ambil token JWT dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("Authorization header is missing")
	}

	// Split token dari header
	tokenString := strings.Split(authHeader, " ")[1]

	// Parse dan verifikasi token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan algoritma token adalah HMAC dan secret key sudah didefinisikan sebelumnya
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return 0, err
	}

	// Pastikan token valid dan sudah terverifikasi
	if !token.Valid {
		return 0, errors.New("Invalid token")
	}

	// Ambil user ID dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to parse token claims")
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("Invalid user ID format in token claims")
	}

	userID := uint(userIDFloat)
	return userID, nil
}
