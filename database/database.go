package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func DBConnection() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error Loading .env file")
	}

	// 환경변수를 이용하여서 DB 접속 정보를 가지고 옴.
	DB_ACCOUNT := os.Getenv("DB_ACCOUNT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", DB_ACCOUNT+":"+DB_PASSWORD+"@tcp("+DB_HOST+")/"+DB_NAME)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	DB = db
}
