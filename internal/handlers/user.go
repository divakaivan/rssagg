package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/utils"
	"github.com/google/uuid"
)

// HandlerCreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param name body string true "Name of the user"
// @Success 201 {object} database.User
// @Failure 400 "{object}" "error"
// @Router /users [post]
func (apiCfg *API) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		utils.RespondWithError(w, 400, fmt.Sprintf("Could not create user: %v", err))
		return
	}

	utils.RespondWithJSON(w, 201, utils.DatabaseUserToUser(user))
}

// HandlerGetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID. Requires authentication.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID of the user"
// @Success 200 {object} database.User
// @Failure 404 "{object}" "error"
// @Router /users/{id} [get]
func (apiCfg *API) HandlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	utils.RespondWithJSON(w, 200, utils.DatabaseUserToUser(user))
}

// HandlerGetPostsForUser godoc
// @Summary Get posts for a user
// @Description Get posts for a user
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit of posts to return"
// @Param offset query int false "Offset of posts to return"
// @Success 200 {array} database.Post
// @Failure 500 "{object}" "error"
// @Router /posts [get]
func (apiCfg *API) HandlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	if r.URL.Query().Get("limit") != "" && r.URL.Query().Get("offset") != "" {
		apiCfg.HandlerGetPostsForUserWithPagination(w, r, user)
		return
	}

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Could not get posts for user: %v", err))
		return
	}
	utils.RespondWithJSON(w, 200, utils.DatabasePostsToPosts(posts))
}

// HandlerGetPostsForUserWithPagination godoc
// @Summary Get posts for a user with pagination
// @Description Get posts for a user with pagination
// @Tags posts
// @Accept json
// @Produce json
// @Param limit query int false "Limit of posts to return"
// @Param offset query int false "Offset of posts to return"
// @Success 200 {array} database.Post
// @Failure 500 "{object}" "error"
// @Router /posts [get]
func (apiCfg *API) HandlerGetPostsForUserWithPagination(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 2
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	offset := 0
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	posts, err := apiCfg.DB.GetPostsForUserWithPagination(r.Context(), database.GetPostsForUserWithPaginationParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Could not get paginated posts for user: %v", err))
		return
	}
	utils.RespondWithJSON(w, 200, utils.DatabasePostsToPosts(posts))
}
