package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/podops/podops/internal/gql/graph/generated"
	"github.com/podops/podops/internal/gql/graph/model"
)

func (r *queryResolver) Show(ctx context.Context, name *string) (*model.Show, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Episode(ctx context.Context, guid *string) (*model.Episode, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Recent(ctx context.Context, max int) ([]*model.Show, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Popular(ctx context.Context, max int) ([]*model.Show, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
