package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/brus-fabrika/rssapi/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerGetFeedByUserId(w http.ResponseWriter, r *http.Request, user database.User) error {
	feeds, err := apiCfg.DB.GetFeedsByUserId(r.Context(), user.ID)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, feeds)

	return nil
}

func (apiCfg *apiConfig) handlerFeed(w http.ResponseWriter, r *http.Request, user database.User) error {

	type feedParameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	var newFeed feedParameters

	json.NewDecoder(r.Body).Decode(&newFeed)

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      newFeed.Name,
		Url:       newFeed.Url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	WriteJson(w, http.StatusCreated, feed)
	return nil
}
