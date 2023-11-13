package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brus-fabrika/rssapi/internal/database"
	"github.com/go-chi/chi"
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

	feed, err := UrlToFeed("https://feeds.nos.nl/jeugdjournaal")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(feed)

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

	router := chi.NewRouter()
	srv := &http.Server{
		Addr:    ":" + apiCfg.Port,
		Handler: router,
	}

	v1Router := chi.NewRouter()
	v1Router.Get("/status", makeHttpHandler(apiCfg.logIt(apiCfg.handlerStatus)))

	v1Router.Get("/users", makeHttpHandler(apiCfg.logIt(apiCfg.handlerGetUsers)))
	v1Router.Get("/user", makeHttpHandler(apiCfg.logIt(apiCfg.middlewareAuth(apiCfg.handlerGetUser))))
	v1Router.Get("/user/:id", makeHttpHandler(apiCfg.logIt(apiCfg.handlerGetUserById)))
	v1Router.Post("/user", makeHttpHandler(apiCfg.logIt(apiCfg.handlerCreateUser)))

	v1Router.Get("/feeds", makeHttpHandler(apiCfg.logIt(apiCfg.middlewareAuth(apiCfg.handlerGetFeedByUserId))))
	v1Router.Post("/feed", makeHttpHandler(apiCfg.logIt(apiCfg.middlewareAuth(apiCfg.handlerFeed))))

	v1Router.Get("/feeds/follow", makeHttpHandler(apiCfg.logIt(apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollowByUserId))))
	v1Router.Post("/feed/follow", makeHttpHandler(apiCfg.logIt(apiCfg.middlewareAuth(apiCfg.handlerFeedFollow))))

	router.Mount("/v1", v1Router)

	// Start the server
	log.Println("Starting server on port", apiCfg.Port)

	startScraping(apiCfg.DB, 1, time.Second*10)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (apiCfg *apiConfig) handlerStatus(w http.ResponseWriter, r *http.Request) error {
	return WriteJson(w, http.StatusOK, struct{}{})
}
