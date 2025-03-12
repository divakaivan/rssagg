package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// HandlerCreateFeedFollow godoc
// @Summary Create a new feed follow
// @Description Create a new feed follow. Requires authentication.
// @Tags feed_follows
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feed_id body string true "Feed ID"
// @Success 201 "{object}" "FeedFollow"
// @Failure 400 "{object}" "ErrorResponse"
// @Router /feed_follows [post]
func (apiCfg *API) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not create feed follow: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, utils.DatabaseFeedFollowToFeedFollow(feedFollow))
}

// HandlerGetFeedFollows godoc
// @Summary Get all feed follows
// @Description Get all feed follows. Requires authentication.
// @Tags feed_follows
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 "{object}" "[]FeedFollow"
// @Failure 400 "{object}" "ErrorResponse"
// @Router /feed_follows [get]
func (apiCfg *API) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not get feed follows: %v", err))
		return
	}

	utils.RespondWithJSON(w, 200, utils.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

// HandlerDeleteFeedFollow godoc
// @Summary Delete a feed follow
// @Description Delete a feed follow. Requires authentication.
// @Tags feed_follows
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedFollowID path string true "Feed Follow ID"
// @Success 200 "{object}" "struct{}"
// @Failure 400 "{object}" "ErrorResponse"
// @Router /feed_follows/{feedFollowID} [delete]
func (apiCfg *API) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not parse feed follow ID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not delete feed follow: %v", err))
		return
	}
	utils.RespondWithJSON(w, 200, struct{}{})
}
