//go:generate go run github.com/vektah/dataloaden -keys int github.com/shufo/go-graphql-boilerplate/models.User

package dataloader

import (
	"context"
	"net/http"
)

type contextKey struct {
	name string
}

var UserLoaderKey = &contextKey{name: "userLoader"}

type loaders struct {
	UserByTodo *UserLoader
}

func DataloaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}
		ldrs.UserByTodo = NewUserLoader(NewUserLoaderConfig(r))

		ctx := context.WithValue(r.Context(), UserLoaderKey, ldrs)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func CtxLoaders(ctx context.Context) loaders {
	return ctx.Value(UserLoaderKey).(loaders)
}
