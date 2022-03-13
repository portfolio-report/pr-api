package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/portfolio-report/pr-api/graph"
	"github.com/portfolio-report/pr-api/graph/generated"
)

func (h *Handler) GraphqlHandler() gin.HandlerFunc {
	graphHandler := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			DB:               h.DB,
			Validate:         h.validate,
			UserService:      h.UserService,
			SessionService:   h.SessionService,
			PortfolioService: h.PortfolioService,
			SecurityService:  h.SecurityService,
		},
	}))

	return func(c *gin.Context) {
		graphHandler.ServeHTTP(c.Writer, c.Request)
	}
}

func (*Handler) PlaygroundHandler(graphqlUrl string) gin.HandlerFunc {
	playgroundHandler := playground.Handler("GraphQL", graphqlUrl)

	return func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	}
}
