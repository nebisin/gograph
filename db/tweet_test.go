package db

import (
	"context"
	"github.com/nebisin/gograph/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTweet(t *testing.T, user User) Tweet {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := CreateTweetParams{
		Content:  util.RandomContent(),
		AuthorID: user.ID,
	}
	tweet, err := testRepository.CreateTweet(ctx, args)
	require.NoError(t, err)
	require.NotEmpty(t, tweet)

	require.Equal(t, tweet.Content, args.Content)
	require.Equal(t, tweet.AuthorId, user.ID)
	require.NotZero(t, tweet.CreatedAt)
	require.NotZero(t, tweet.UpdatedAt)

	return tweet
}

func TestRepository_CreateTweet(t *testing.T) {
	user := createRandomUser(t)

	createRandomTweet(t, user)
}

func TestRepository_GetTweet(t *testing.T) {
	user := createRandomUser(t)

	tweet1 := createRandomTweet(t, user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tweet2, err := testRepository.GetTweet(ctx, tweet1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tweet2)

	require.Equal(t, tweet2.ID, tweet1.ID)
	require.Equal(t, tweet2.Content, tweet1.Content)
	require.Equal(t, tweet2.AuthorId, tweet1.AuthorId)
	require.WithinDuration(t, tweet2.CreatedAt, tweet1.CreatedAt, time.Second)
	require.WithinDuration(t, tweet2.UpdatedAt, tweet1.UpdatedAt, time.Second)
}

func TestRepository_ListTweet(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomTweet(t, user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tweets, err := testRepository.ListTweet(ctx, 5, 1)
	require.NoError(t, err)
	require.Len(t, tweets, 5)

	for _, tweet := range tweets {
		require.NotEmpty(t, tweet)
		require.Equal(t, tweet.AuthorId, user.ID)
	}
}

func TestRepository_DeleteTweet(t *testing.T) {
	user := createRandomUser(t)

	tweet1 := createRandomTweet(t, user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := testRepository.DeleteTweet(ctx, tweet1.ID)
	require.NoError(t, err)

	tweet2, err := testRepository.GetTweet(ctx, tweet1.ID)
	require.Error(t, err)
	require.Empty(t, tweet2)
}

func TestRepository_UpdateTweet(t *testing.T) {
	user := createRandomUser(t)
	tweet1 := createRandomTweet(t, user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := UpdateTweetParams{
		ID:      tweet1.ID,
		Content: util.RandomContent(),
	}

	tweet2, err := testRepository.UpdateTweet(ctx, args)
	require.NoError(t, err)
	require.NotEmpty(t, tweet2)

	require.Equal(t, tweet2.ID, tweet1.ID)
	require.Equal(t, tweet2.AuthorId, tweet1.AuthorId)
	require.Equal(t, tweet2.Content, args.Content)
	require.WithinDuration(t, tweet2.CreatedAt, tweet1.CreatedAt, time.Second)
	require.NotZero(t, tweet2.UpdatedAt)
}

func TestRepository_ListTweetByAuthor(t *testing.T) {
	user := createRandomUser(t)

	for i := 0; i < 5; i++ {
		createRandomTweet(t, user)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tweets, err := testRepository.ListTweetByAuthor(ctx, user.ID, 5, 1)
	require.NoError(t, err)
	require.Len(t, tweets, 5)

	for _, tweet := range tweets {
		require.NotEmpty(t, tweet)
		require.Equal(t, tweet.AuthorId, user.ID)
	}
}
