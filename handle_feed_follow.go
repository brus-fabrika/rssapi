package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/brus-fabrika/rssapi/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) error {
	type feedParameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	var newFeed feedParameters

	json.NewDecoder(r.Body).Decode(&newFeed)

	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    newFeed.FeedId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	WriteJson(w, http.StatusCreated, feed)
	return nil
}

func (apiCfg *apiConfig) handlerGetFeedFollowByUserId(w http.ResponseWriter, r *http.Request, user database.User) error {
	feeds_follow, err := apiCfg.DB.GetFeedsFollowByUserId(r.Context(), user.ID)
	if err != nil {
		return err
	}

	WriteJson(w, http.StatusOK, feeds_follow)

	return nil
}
