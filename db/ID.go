package db

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

func MarshalID(id primitive.ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(fmt.Sprintf(`%q`, id.Hex())))
	})
}

func UnmarshalID(v interface{}) (primitive.ObjectID, error) {
	switch v := v.(type) {
	case string:
		return primitive.ObjectIDFromHex(v)
	default:
		return primitive.ObjectID{}, fmt.Errorf("%T is not an id", v)
	}
}
