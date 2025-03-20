package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)


var DB *sql.DB

func ConnectDB (){
	

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")


	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Cek koneksi ke database
	if err := DB.Ping(); err != nil {
		log.Fatal("Tidak dapat menjangkau database:", err)
	}

	fmt.Println("Berhasil terhubung ke database")
}