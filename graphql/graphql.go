package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	"github.com/podops/podops/apiv1"
	"github.com/podops/podops/graphql/graph"
	"github.com/podops/podops/graphql/graph/generated"
)

// GraphqlEndpoint maps the Graphql handler to gin
func GraphqlEndpoint() echo.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.CreateResolver()}))

	return func(e echo.Context) error {
		h.ServeHTTP(e.Response(), e.Request())
		return nil
	}
}

// GraphqlPlaygroundEndpoint maps the Playground handler to gin
func GraphqlPlaygroundEndpoint() echo.HandlerFunc {
	h := playground.Handler("GraphQL", apiv1.GraphqlNamespacePrefix+apiv1.GraphqlRoute)

	return func(e echo.Context) error {
		h.ServeHTTP(e.Response(), e.Request())
		return nil
	}
}
