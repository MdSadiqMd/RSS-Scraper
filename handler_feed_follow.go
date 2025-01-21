package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MdSadiqMd/RSS-Scraper/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type Parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
	}

	follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, "Unable to create feed: "+err.Error())
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(follow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, "Unable to get feed follows: "+err.Error())
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(follows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feed_id")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, "Invalid feed follow id")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, "Unable to delete feed follow: "+err.Error())
		return
	}

	respondWithJSON(w, 200, map[string]string{"status": "deleted"})
}
