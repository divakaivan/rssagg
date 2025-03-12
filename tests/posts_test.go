package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/divakaivan/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	user := CreateRandomUser(t)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.ID)

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test Feed123",
		Url:       "https://example123.com",
		UserID:    user.ID,
	}
	feed, _ := testQueries.CreateFeed(ctx, feedParams)

	params := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       "Test Posasdasdat",
		Url:         "https://example_posthahahahaha.com",
		Content:     "This is a test post",
		PublishedAt: time.Now().UTC().Format(time.RFC3339),
		FeedID:      feed.ID,
	}

	post, err := testQueries.CreatePost(ctx, params)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	assert.Equal(t, params.Url, post.Url)
	assert.Equal(t, params.Title, post.Title)
}

func TestPostRepository_GetPostsByUser(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	name := "Random Test User2"
	user, _ := testQueries.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      "Test Feed1234",
		Url:       "https://example1234.com",
		UserID:    user.ID,
	}
	feed, _ := testQueries.CreateFeed(ctx, feedParams)

	var createdPosts []database.Post

	for i := 0; i < 11; i++ {
		post, err := testQueries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       fmt.Sprintf("Test Post %d", i),
			Summary:     "Test summary",
			Content:     "This is a test post",
			Url:         fmt.Sprintf("https://example.com/%d", i),
			PublishedAt: time.Now().UTC().Format(time.RFC3339),
			FeedID:      feed.ID,
		})
		require.NoError(t, err)
		require.NotEmpty(t, post)
		require.NotEmpty(t, post.ID)
		createdPosts = append(createdPosts, post)
	}

	posts, err := testQueries.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: feed.UserID,
		Limit:  10,
	})

	require.NoError(t, err)
	require.Len(t, posts, 10)
	for i, post := range posts {
		assert.Equal(t, createdPosts[i].ID, post.ID)
		assert.Equal(t, feed.ID, post.FeedID)
	}
}
