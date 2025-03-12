package handlers

import (
	"net/http"

	"github.com/divakaivan/rssagg/internal/utils"
)

// HandlerReadiness godoc
// @Summary Check if the service is healthy
// @Description Check if the service is healthy
// @Tags healthz
// @Produce json
// @Router /healthz [get]
func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, struct{}{})
}
