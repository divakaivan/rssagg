package tests

import (
	"context"
	"testing"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) database.User {
	t.Helper()

	name := "Random Test User"
	user, err := testQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotEmpty(t, user.ID)
	require.Equal(t, name, user.Name)
	require.NotEmpty(t, user.CreatedAt)
	require.NotEmpty(t, user.UpdatedAt)
	require.NotEmpty(t, user.ApiKey)

	return user
}

func TestQueries_CreateUser(t *testing.T) {
	user := CreateRandomUser(t)

	testDB.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)

	user2, err := testQueries.GetUserByAPIKey(context.Background(), user.ApiKey)
	require.Error(t, err)
	assert.Empty(t, user2)
}

func TestQueries_GetUserFromApiKey(t *testing.T) {
	user := CreateRandomUser(t)

	user2, err := testQueries.GetUserByAPIKey(context.Background(), user.ApiKey)
	require.NoError(t, err)
	assert.Equal(t, user, user2)

	testDB.QueryContext(context.Background(), "Delete from users where id = $1", user.ID)

	user3, err := testQueries.GetUserByAPIKey(context.Background(), user.ApiKey)
	require.Error(t, err)
	assert.Empty(t, user3)
}
