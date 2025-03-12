package handlers

import (
	"net/http"

	"github.com/divakaivan/rssagg/internal/utils"
)

// HandlerErr godoc
// @Summary Something went wrong
// @Description Something went wrong
// @Tags err
// @Produce json
// @Router /err [get]
func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 400, "Something went wrong")
}
