package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/divakaivan/rssagg/internal/auth"
	"github.com/divakaivan/rssagg/internal/database"
	"github.com/google/uuid"
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
			apiKey = ""
		}
		user, err := api.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			// requires the creation of a default user in the db
			emptyUserID, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
			user.ID = emptyUserID
		}
		err = api.DB.CreateLog(r.Context(), database.CreateLogParams{
			Timestamp:    time.Now(),
			CallerUserID: user.ID,
			RemoteAddr:   r.RemoteAddr,
			Method:       r.Method,
			Url:          r.URL.Path,
			Status:       strconv.Itoa(wrappedWriter.statusCode),
			DurationMs:   duration.Milliseconds(),
		})
		if err != nil {
			log.Print("[Error] Cannot log request", err)
		}
	})
}
