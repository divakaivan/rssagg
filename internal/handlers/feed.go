package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/utils"
	"github.com/google/uuid"
)

// HandlerCreateFeed godoc
// @Summary Create a new feed
// @Description Create a new feed. Requires authentication.
// @Tags feeds
// @Accept json
// @Param name body string true "Name of the feed"
// @Param url body string true "URL of the feed"
// @Produce json
// @Security ApiKeyAuth
// @Success 201 "Feed"
// @Failure 400 "ErrorResponse"
// @Router /feeds [post]
func (apiCfg *API) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not create feed: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, utils.DatabaseFeedToFeed(feed))
}

// HandlerGetFeeds godoc
// @Summary Get all feeds
// @Description Get all feeds
// @Tags feeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 "[]Feed"
// @Failure 400 "ErrorResponse"
// @Router /feeds [get]
func (apiCfg *API) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not get feeds: %v", err))
		return
	}

	utils.RespondWithJSON(w, 200, utils.DatabaseFeedsToFeeds(feeds))
}
