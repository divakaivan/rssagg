package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/divakaivan/rssagg/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHandlerErr(t *testing.T) {
	req, err := http.NewRequest("GET", "/some-endpoint", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandlerErr)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error": "Something went wrong"}`, rr.Body.String())
}
