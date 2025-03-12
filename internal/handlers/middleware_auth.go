package handlers

import (
	"fmt"
	"net/http"

	"github.com/divakaivan/rssagg/internal/auth"
	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg API) AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("Could not get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}
