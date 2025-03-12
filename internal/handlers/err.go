package handlers

import (
	"net/http"

	"github.com/divakaivan/rssagg/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}
