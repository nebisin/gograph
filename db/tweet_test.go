package db

import (
	"context"
	"github.com/nebisin/gograph/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	t.Run("valid content", func(t *testing.T) {
		createRandomTweet(t, user)
	})

	t.Run("invalid content", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		args := CreateTweetParams{
			Content:  util.RandomString(190),
			AuthorID: user.ID,
		}
		tweet, err := testRepository.CreateTweet(ctx, args)
		require.Error(t, err)
		require.Empty(t, tweet)
	})
}

func TestRepository_GetTweet(t *testing.T) {
	user := createRandomUser(t)

	tweet1 := createRandomTweet(t, user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with existed id", func(t *testing.T) {
		tweet2, err := testRepository.GetTweet(ctx, tweet1.ID)
		require.NoError(t, err)
		require.NotEmpty(t, tweet2)

		require.Equal(t, tweet2.ID, tweet1.ID)
		require.Equal(t, tweet2.Content, tweet1.Content)
		require.Equal(t, tweet2.AuthorId, tweet1.AuthorId)
		require.WithinDuration(t, tweet2.CreatedAt, tweet1.CreatedAt, time.Second)
		require.WithinDuration(t, tweet2.UpdatedAt, tweet1.UpdatedAt, time.Second)
	})

	t.Run("with not existed id", func(t *testing.T) {
		tweet2, err := testRepository.GetTweet(ctx, primitive.NewObjectID())
		require.Error(t, err)
		require.Empty(t, tweet2)
	})
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

	t.Run("with existed id", func(t *testing.T) {
		err := testRepository.DeleteTweet(ctx, tweet1.ID)
		require.NoError(t, err)

		tweet2, err := testRepository.GetTweet(ctx, tweet1.ID)
		require.Error(t, err)
		require.Empty(t, tweet2)
	})

	t.Run("with not existed id", func(t *testing.T) {
		err := testRepository.DeleteTweet(ctx, tweet1.ID)
		require.Error(t, err)
	})
}

func TestRepository_UpdateTweet(t *testing.T) {
	user := createRandomUser(t)
	tweet1 := createRandomTweet(t, user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with valid params", func(t *testing.T) {
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
	})

	t.Run("with invalid content", func(t *testing.T) {
		args := UpdateTweetParams{
			ID:      tweet1.ID,
			Content: util.RandomString(190),
		}

		tweet2, err := testRepository.UpdateTweet(ctx, args)
		require.Error(t, err)
		require.Empty(t, tweet2)
	})

	t.Run("with not existed id", func(t *testing.T) {
		args := UpdateTweetParams{
			ID:      primitive.NewObjectID(),
			Content: util.RandomContent(),
		}

		tweet2, err := testRepository.UpdateTweet(ctx, args)
		require.Error(t, err)
		require.Empty(t, tweet2)
	})

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
