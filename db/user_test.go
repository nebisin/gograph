package db

import (
	"context"
	"github.com/nebisin/gograph/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := RegisterParams{
		Email:       util.RandomEmail(),
		Password:    util.RandomPassword(),
		DisplayName: util.RandomDisplayName(),
	}

	user, err := testRepository.CreateUser(ctx, args)
	require.NoError(t, err)

	require.Equal(t, user.Email, args.Email)
	require.Equal(t, user.DisplayName, args.DisplayName)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	err = util.CheckPassword(args.Password, user.Password)
	require.NoError(t, err)

	return user
}

func TestRepository_CreateUser(t *testing.T) {
	user := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with invalid email", func(t *testing.T) {
		args := RegisterParams{
			Email:       "someinvalidemail.com",
			Password:    util.RandomPassword(),
			DisplayName: util.RandomDisplayName(),
		}

		user, err := testRepository.CreateUser(ctx, args)
		require.Error(t, err)
		require.Empty(t, user)
	})

	t.Run("with taken email", func(t *testing.T) {
		args := RegisterParams{
			Email:       user.Email,
			Password:    util.RandomPassword(),
			DisplayName: util.RandomDisplayName(),
		}

		user, err := testRepository.CreateUser(ctx, args)
		require.Error(t, err)
		require.Empty(t, user)
	})
}

func TestRepository_GetUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with valid id", func(t *testing.T) {
		user2, err := testRepository.GetUser(ctx, user1.ID)
		require.NoError(t, err)
		require.NotEmpty(t, user2)

		require.Equal(t, user2.ID, user1.ID)
		require.Equal(t, user2.Email, user1.Email)
		require.Equal(t, user2.Password, user1.Password)
		require.Equal(t, user2.DisplayName, user1.DisplayName)

		require.WithinDuration(t, user2.CreatedAt, user1.CreatedAt, time.Second)
		require.WithinDuration(t, user2.UpdatedAt, user1.UpdatedAt, time.Second)
	})

	t.Run("with not existed id", func(t *testing.T) {
		user2, err := testRepository.GetUser(ctx, primitive.NewObjectID())
		require.Error(t, err)
		require.Empty(t, user2)
	})

}

func TestRepository_GetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with existed email", func(t *testing.T) {
		user2, err := testRepository.GetUserByEmail(ctx, user1.Email)
		require.NoError(t, err)
		require.NotEmpty(t, user2)

		require.Equal(t, user2.Email, user1.Email)
		require.Equal(t, user2.Password, user1.Password)
		require.Equal(t, user2.DisplayName, user1.DisplayName)

		require.WithinDuration(t, user2.CreatedAt, user1.CreatedAt, time.Second)
		require.WithinDuration(t, user2.UpdatedAt, user1.UpdatedAt, time.Second)
	})

	t.Run("with not existed email", func(t *testing.T) {
		user2, err := testRepository.GetUserByEmail(ctx, util.RandomEmail())
		require.Error(t, err)
		require.Empty(t, user2)
	})

}

func TestRepository_UpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with solid params", func(t *testing.T) {
		params := UpdateUserParams{
			DisplayName: util.RandomDisplayName(),
			Email: util.RandomEmail(),
			Password: util.RandomPassword(),
		}

		user2, err := testRepository.UpdateUser(ctx, user1.ID, params)
		require.NoError(t, err)
		require.NotEmpty(t, user2)

		require.Equal(t, user1.ID, user2.ID)
		require.Equal(t, user2.Email, params.Email)
		require.Equal(t, user2.DisplayName, params.DisplayName)
		require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

		err = util.CheckPassword(params.Password, user2.Password)
		require.NoError(t, err)
	})

	t.Run("with invalid password", func(t *testing.T) {
		params := UpdateUserParams{
			Password: "123",
		}

		user2, err := testRepository.UpdateUser(ctx, user1.ID, params)
		require.Error(t, err)
		require.Empty(t, user2)
	})

	t.Run("with invalid email", func(t *testing.T) {
		params := UpdateUserParams{
			Email: "invalidemail.com",
		}

		user2, err := testRepository.UpdateUser(ctx, user1.ID, params)
		require.Error(t, err)
		require.Empty(t, user2)
	})

	t.Run("with taken email", func(t *testing.T) {
		user3 := createRandomUser(t)

		params := UpdateUserParams{
			Email: user3.Email,
		}

		user2, err := testRepository.UpdateUser(ctx, user1.ID, params)
		require.Error(t, err)
		require.Empty(t, user2)
	})

	t.Run("with invalid display name", func(t *testing.T) {
		params := UpdateUserParams{
			DisplayName: ".:.",
		}

		user2, err := testRepository.UpdateUser(ctx, user1.ID, params)
		require.Error(t, err)
		require.Empty(t, user2)
	})
}

func TestRepository_DeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("with existed user", func(t *testing.T) {
		err := testRepository.DeleteUser(ctx, user1.ID)
		require.NoError(t, err)

		user2, err := testRepository.GetUser(ctx, user1.ID)
		require.Error(t, err)
		require.Empty(t, user2)
	})

	t.Run("with non existed user", func(t *testing.T) {
		err := testRepository.DeleteUser(ctx, user1.ID)
		require.Error(t, err)
	})
}