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

type ApiFunc func(w http.ResponseWriter, r *http.Request) error
type ApiError struct {
	Error string `json:"error"`
}

func WriteJson(w http.ResponseWriter, httpStatus int, v any) error {
	w.WriteHeader(httpStatus)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandler(fn ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type apiConfig struct {
	Port string
	DB   *database.Queries
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
		Port: portString,
		DB:   database.New(db),
	}

	apiCfg.Run()
}

func (apiCfg *apiConfig) Run() {
	http.HandleFunc("/users/create", makeHttpHandler(apiCfg.handlerCreateUser))
	http.HandleFunc("/users", makeHttpHandler(apiCfg.handlerGetUsers))

	// Start the server
	log.Println("Starting server on port", apiCfg.Port)
	http.ListenAndServe(":"+apiCfg.Port, nil)
}

func (apiCfg *apiConfig) handlerGetUserById(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (apiCfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) error {

	log.Println("Getting all users")

	users, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, users)

	return nil
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) error {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return err
	}

	log.Println("Creating user with name:", params.Name)

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, user)
	return nil
}
