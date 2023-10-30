package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// docker run --rm -v "$((Get-Item .).FullName):/src" -w /src sqlc/sqlc generate

const (
	DB_USER     = "root"
	DB_PASSWORD = "root"
	DB_NAME     = "gotest"
)

func main() {

	// Get the PORT environment variable
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment is not set")
	}

	log.Println("PORT is set to", portString)

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("Database connection established")

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database table created")
}
