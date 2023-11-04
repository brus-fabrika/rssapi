package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/brus-fabrika/rssapi/internal/auth"
	"github.com/brus-fabrika/rssapi/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetUserById(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) error {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		return err
	}

	user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, user)

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

	WriteJson(w, http.StatusCreated, user)
	return nil
}
