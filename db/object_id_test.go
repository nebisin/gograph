package db

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUnmarshalID(t *testing.T) {
	t.Run("marshal valid id", func(t *testing.T) {
		objectId1 := primitive.NewObjectID()
		stringId := objectId1.Hex()

		objectId2, err := UnmarshalID(stringId)
		require.NoError(t, err)
		require.Equal(t, objectId2, objectId1)
	})

	t.Run("marshal invalid id", func(t *testing.T) {
		objectId2, err := UnmarshalID(123)
		require.Error(t, err)
		require.Zero(t, objectId2)
	})

}
