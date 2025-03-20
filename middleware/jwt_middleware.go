package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Initialize JWT Key
var jwtKey []byte
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET")) 
}

// Define a custom type for context keys
type contextKey string

const ClaimsKey contextKey = "claims"

// Middleware for checking JWT Token
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Println("No Token Provided!")
			http.Error(w, "No Token Provided!", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			fmt.Println("Invalid Token Format!")
			http.Error(w, "Invalid Token Format!", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// fmt.Println("Token: ", token)
		// fmt.Println("Claims: ", claims)

		if err != nil || !token.Valid {
			fmt.Println("Invalid Token!")
			http.Error(w, "Invalid Token!", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
