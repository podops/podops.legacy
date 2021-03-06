package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/fupas/commons/pkg/util"
	ds "github.com/fupas/platform/pkg/platform"
	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/internal/gql/graph/generated"
	"github.com/podops/podops/internal/gql/graph/model"
	"github.com/podops/podops/internal/platform"
	"github.com/podops/podops/pkg/backend"
)

func (r *queryResolver) Show(ctx context.Context, name *string) (*model.Show, error) {
	if r.ShowLoader == nil {
		log.Fatal("panic: missing show loader")
	}

	data, err := r.ShowLoader.Load(ctx, *name)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}
	show := data.(*model.Show)

	// list all episodes, excluding future (i.e. unpublished) ones, descending order
	var er []*a.Resource
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(backend.DatastoreResources).Filter("ParentGUID =", show.GUID).Filter("Kind =", a.ResourceEpisode).Filter("Published <", util.Timestamp()).Order("-Published"), &er); err != nil {
		platform.ReportError(err)
		return nil, err
	}

	if er != nil {
		episodes := make([]*model.Episode, len(er))
		for i := range er {
			e, err := r.Episode(ctx, &er[i].GUID)
			if err != nil {
				// no need to log, already done in r.Episode
				return nil, err
			}
			episodes[i] = e
		}
		show.Episodes = episodes
	}

	return show, nil
}

func (r *queryResolver) Episode(ctx context.Context, guid *string) (*model.Episode, error) {
	if r.EpisodeLoader == nil {
		log.Fatal("panic: missing episode loader")
	}

	data, err := r.EpisodeLoader.Load(ctx, *guid)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}
	return data.(*model.Episode), nil
}

func (r *queryResolver) Recent(ctx context.Context, max int) ([]*model.Show, error) {
	var sh []*a.Production
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(backend.DatastoreProductions).Filter("BuildDate >", 0).Order("-BuildDate").Limit(max), &sh); err != nil {
		platform.ReportError(err)
		return nil, err
	}

	var shows []*model.Show
	if sh != nil {
		shows = make([]*model.Show, len(sh))
		for i := range sh {
			show, err := r.ShowLoader.Load(ctx, sh[i].Name)
			if err != nil {
				platform.ReportError(err)
				return nil, err
			}
			shows[i] = show.(*model.Show)
		}
	}

	return shows, nil
}

func (r *queryResolver) Popular(ctx context.Context, max int) ([]*model.Show, error) {
	var sh []*a.Production
	// FIXME change this once we have metrics on show subscriptions
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(backend.DatastoreProductions).Filter("BuildDate >", 0).Order("-BuildDate").Limit(max), &sh); err != nil {
		platform.ReportError(err)
		return nil, err
	}

	var shows []*model.Show
	if sh != nil {
		shows = make([]*model.Show, len(sh))
		for i := range sh {
			show, err := r.ShowLoader.Load(ctx, sh[i].Name)
			if err != nil {
				platform.ReportError(err)
				return nil, err
			}
			shows[i] = show.(*model.Show)
		}
	}

	return shows, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
