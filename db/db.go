package db

import (
	"database/sql"
	"os"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func InitializeDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error Loading  .env file")
	}
	dbdriver := os.Getenv("DBDRIVER")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	database := os.Getenv("DATABASE")
	PORT := os.Getenv("PORT")

	DBURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, PORT, database)
	db, err = sql.Open(dbdriver, DBURL)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Success connected to database.")
}

func GetDB() *sql.DB {
	return db
}
