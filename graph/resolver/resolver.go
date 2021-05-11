package resolver

import (
	"github.com/nebisin/gograph/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository *db.Repository
}
