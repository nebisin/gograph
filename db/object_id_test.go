package db

import (
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUnmarshalID(t *testing.T) {
	objectId1 := primitive.NewObjectID()
	stringId := objectId1.Hex()

	objectId2, err := UnmarshalID(stringId)
	require.NoError(t, err)
	require.Equal(t, objectId2, objectId1)
}
