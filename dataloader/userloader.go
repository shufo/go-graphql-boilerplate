package dataloader

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/shufo/go-graphql-boilerplate/models"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func NewUserLoaderConfig(r *http.Request) UserLoaderConfig {
	return UserLoaderConfig{
		MaxBatch: 100,
		Wait:     1 * time.Millisecond,
		Fetch: func(ids []int) ([]*models.User, []error) {
			db := r.Context().Value("db").(*sql.DB)

			ctx := context.Background()

			s := make([]interface{}, len(ids))

			for i, v := range ids {
				s[i] = v
			}

			users, _ := models.Users(qm.OrIn("id in ?", s...)).All(ctx, db)

			results := make([]*models.User, len(users))

			for i, key := range ids {
				for _, user := range users {
					if key == int(user.ID) {
						results[i] = user
					}
				}
			}

			return results, nil
		},
	}
}
