package auth

import (
	"context"
	"net/http"

	"github.com/shufo/go-graphql-boilerplate/logger"

	"github.com/go-chi/jwtauth"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses

var UserCtxKey = &ContextKey{name: "user"}

type ContextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get claims from current context
		_, claims, _ := jwtauth.FromContext(r.Context())

		// output user id to log if it exists
		if claims["user_id"] != nil {
			logger.LogEntrySetField(r, "user_id", claims["user_id"])
		}

		// put claims in context
		ctx := context.WithValue(r.Context(), UserCtxKey, claims)

		// and call the next with our new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) jwtauth.Claims {
	raw, _ := ctx.Value(UserCtxKey).(jwtauth.Claims)
	return raw
}
