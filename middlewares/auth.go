package middlewares

import (
	"context"
	"github.com/nebisin/gograph/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeaderKey)
			if len(authHeader) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			fields := strings.Fields(authHeader)
			if len(fields) < 2 {
				http.Error(w,"invalid authorization header", http.StatusUnauthorized)
				return
			}

			authType := strings.ToLower(fields[0])
			if authType != authorizationTypeBearer {
				http.Error(w,"unsupported authorization type", http.StatusUnauthorized)
				return
			}

			accessToken := fields[1]

			payload, err := token.VerifyToken(accessToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authorizationPayloadKey, payload)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func AuthContext(ctx context.Context) *token.Payload {
	raw, _ := ctx.Value(authorizationPayloadKey).(*token.Payload)
	return raw
}