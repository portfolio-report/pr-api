package dataloaders

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph/model"
	"github.com/portfolio-report/pr-api/libs/dataloader"
)

type contextKey string

const dataloaderCtxKey = contextKey("dataloaders")

// Loaders holds references to dataloaders
type Loaders struct {
	UserByID *dataloader.Dataloader[int, *model.User]
}

func newLoaders(ctx context.Context, userService model.UserService) *Loaders {
	return &Loaders{
		UserByID: dataloader.New(
			func(keys []int) ([]*model.User, []error) {
				users, _ := userService.GetByIDs(keys)

				// map by id
				usersByID := make(map[int]*model.User, len(users))
				for _, u := range users {
					usersByID[u.ID] = u
				}
				// list in order of keys
				result := make([]*model.User, len(keys))
				for i, key := range keys {
					result[i] = usersByID[key]
				}

				return result, nil
			}, 1*time.Millisecond, 0),
	}
}

// Middleware returns a middleware function that attaches loaders to request context
func Middleware(userService model.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		loaders := newLoaders(ctx, userService)
		ctx = context.WithValue(ctx, dataloaderCtxKey, loaders)
		c.Request = c.Request.WithContext(ctx)
	}
}

// For gets loaders from context
func For(ctx context.Context) *Loaders {
	return ctx.Value(dataloaderCtxKey).(*Loaders)
}
