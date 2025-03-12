package handlers

import "github.com/divakaivan/rssagg/internal/database"

type API struct {
	DB *database.Queries
}

func NewAPI(db *database.Queries) *API {
	return &API{DB: db}
}
