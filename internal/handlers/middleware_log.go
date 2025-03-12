package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/divakaivan/rssagg/internal/auth"
	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/utils"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (api *API) LogToDbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		wrappedWriter := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(startTime)

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := api.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, 400, fmt.Sprintf("Could not get user: %v", err))
			return
		}

		api.DB.CreateLog(r.Context(), database.CreateLogParams{
			Timestamp:    time.Now(),
			CallerUserID: user.ID,
			Method:       r.Method,
			Url:          r.URL.Path,
			Status:       strconv.Itoa(wrappedWriter.statusCode),
			DurationMs:   duration.Milliseconds(),
		})
	})
}
