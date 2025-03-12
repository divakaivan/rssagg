package tests

import (
	"testing"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/divakaivan/rssagg/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestNewAPI(t *testing.T) {
	mockDB := &database.Queries{}
	api := handlers.NewAPI(mockDB)

	assert.NotNil(t, api)
	assert.Equal(t, mockDB, api.DB)
}
