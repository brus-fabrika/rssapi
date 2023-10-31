package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brus-fabrika/rssapi/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// docker run --rm -v "$((Get-Item .).FullName):/src" -w /src sqlc/sqlc generate

type apiConfig struct {
	DB *database.Queries
}

func main() {

	// Get the PORT environment variable
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment is not set")
	}

	dbConnString := os.Getenv("DB_URL")
	if dbConnString == "" {
		log.Fatal("Postgress connection string is not set")
	}

	log.Println("PORT is set to", portString)

	db, err := sql.Open("postgres", dbConnString)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("Database connection established")

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	// Start the server
	log.Println("Starting server on port", portString)
	http.HandleFunc("/users/create", apiCfg.handlerCreateUser)

	// start the server
	http.ListenAndServe(":"+portString, nil)

}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Creating user with name:", params.Name)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
