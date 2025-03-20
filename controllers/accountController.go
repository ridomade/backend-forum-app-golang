package controllers

import (
	"database/sql"
	"encoding/json"
	"firstproject/config"
	"firstproject/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Creating New Account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	
	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Println("Error Reading Data from Request : \n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Reading Data from Request",
		})
		return
	}
	
	// Check if email already registered
	var existingEmail string
	queryCheck := "SELECT email FROM accounts WHERE email = ?"
	err = config.DB.QueryRow(queryCheck, account.Email).Scan(&existingEmail)
	if err == nil {
		fmt.Printf("Email {%s} Already Registerd!", account.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Email Already Registerd!",
		})
		return
	}

	// Hashsing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), 10)
	if err != nil {
		log.Fatal(err)
	}

	// Insert Data to Database
	query	:= "INSERT INTO accounts (email,password) VALUES (?, ?)"
	result, err := config.DB.Exec(query, account.Email, hashedPassword)
	if err != nil {
		fmt.Println("Error Inserting New Data : \n",err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Inserting New Data",
		})
		return
	}

	// Get ID from last inserted data
	id, _ := result.LastInsertId()
	account.ID = int(id)

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// Initialize JWT Key
var jwtKey []byte
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET")) 
}


func LoginAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	var storedPasswordHash string

	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		fmt.Println("Error Reading Request : \n",err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Reading Request",
		})
		return
	}

	// Get account data from database
	query := "SELECT id, password FROM accounts WHERE email = ?"
	err = config.DB.QueryRow(query, account.Email).Scan(&account.ID, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Email Not Found! : \n",err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Email Not Found!",
			})
		} else {	
			fmt.Println("Error when Trying to Get the Data : \n",err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Error when Trying to Get the Data",
			})
		}
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(account.Password))
	if err != nil {
		fmt.Println("Wrong Passowrd! : \n" ,err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Wrong Password!",
		})
		return
	}

	// Create JWT Token
	expirationTime := time.Now().Add(24 * time.Hour) 
	claims := &jwt.RegisteredClaims{
		Subject:   account.Email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Error Creating Token : \n",err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error Creating Token",
		})
		return
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login berhasil!",
		"data":    account,
		"token":   tokenString,	
	})
}
