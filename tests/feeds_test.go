package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomFeed(t *testing.T) database.Feed {
	t.Helper()

	user := CreateRandomUser(t)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.ID)

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test Feed",
		Url:       "https://example.com",
		UserID:    user.ID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, feedParams)
	require.NoError(t, err)
	require.NotEmpty(t, feed)
	assert.Equal(t, feedParams.Name, feed.Name)
	assert.Equal(t, feedParams.Url, feed.Url)
	assert.Equal(t, feedParams.UserID, feed.UserID)
	require.NotEmpty(t, feed.CreatedAt)
	require.NotEmpty(t, feed.UpdatedAt)

	return feed
}

func TestQueries_CreateFeed(t *testing.T) {
	feed := CreateRandomFeed(t)
	assert.NotEmpty(t, feed)
}

func TestQueries_CreateFeed_ViolateFK(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test Feed2",
		Url:       "https://example2.com",
		UserID:    uuid.New(),
	})
	assert.Error(t, err)
	assert.Empty(t, feed)
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	assert.Equal(t, "foreign_key_violation", pqErr.Code.Name())
}

func TestQueries_ListFeeds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := testQueries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test User 123",
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)

	feed, err := testQueries.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test Feed9",
		Url:       "https://example9.com",
		UserID:    user.ID,
	})
	require.NoError(t, err)
	feeds, err := testQueries.GetFeeds(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, feeds)
	assert.Equal(t, 2, len(feeds))
	assert.Equal(t, feed.Name, feeds[1].Name)
	assert.Equal(t, feed.Url, feeds[1].Url)
}
