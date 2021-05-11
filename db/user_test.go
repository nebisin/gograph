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
	require.NotZero(t, user.CreatedAt)

	err = util.CheckPassword(args.Password, user.Password)
	require.NoError(t, err)

	return user
}

func TestRepository_CreateUser(t *testing.T) {
	createRandomUser(t)
}
