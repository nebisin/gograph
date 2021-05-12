package db

import (
	"context"
	"github.com/nebisin/gograph/util"
	"github.com/stretchr/testify/require"
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
	createRandomUser(t)
}

func TestRepository_GetUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user2, err := testRepository.GetUser(ctx, user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.ID, user1.ID)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.Password, user1.Password)
	require.Equal(t, user2.DisplayName, user1.DisplayName)

	require.WithinDuration(t, user2.CreatedAt, user1.CreatedAt, time.Second)
	require.WithinDuration(t, user2.UpdatedAt, user1.UpdatedAt, time.Second)
}

func TestRepository_GetUserByEmail(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user2, err := testRepository.GetUserByEmail(ctx, user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.Password, user1.Password)
	require.Equal(t, user2.DisplayName, user1.DisplayName)

	require.WithinDuration(t, user2.CreatedAt, user1.CreatedAt, time.Second)
	require.WithinDuration(t, user2.UpdatedAt, user1.UpdatedAt, time.Second)
}

func TestRepository_UpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := UpdateUserParams{
		ID: user1.ID,
		DisplayName: util.RandomDisplayName(),
	}

	user2, err := testRepository.UpdateUser(ctx, params)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user2.Email, user1.Email)
	require.Equal(t, user2.DisplayName, params.DisplayName)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}