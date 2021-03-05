package api

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops/internal/gql/graph"
	"github.com/podops/podops/internal/gql/graph/generated"
	"github.com/podops/podops/pkg/api"
)

// GetGraphqlEndpoint maps the Graphql handler to gin
func GetGraphqlEndpoint() echo.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.CreateResolver()}))

	return func(e echo.Context) error {
		h.ServeHTTP(e.Response(), e.Request())
		return nil
	}
}

// GetGraphqlPlaygroundEndpoint maps the Playground handler to gin
func GetGraphqlPlaygroundEndpoint() echo.HandlerFunc {
	h := playground.Handler("GraphQL", api.GraphqlNamespacePrefix+api.GraphqlRoute)

	return func(e echo.Context) error {
		h.ServeHTTP(e.Response(), e.Request())
		return nil
	}
}
