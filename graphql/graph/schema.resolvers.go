package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"

	"github.com/txsvc/platform/v2"
	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/graphql/graph/generated"
	"github.com/podops/podops/graphql/graph/model"
)

func (r *queryResolver) Show(ctx context.Context, name *string, limit int) (*model.Show, error) {

	now := timestamp.Now()

	data, err := r.ShowLoader.Load(ctx, *name)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}
	show := data.(*model.Show)

	// verify that the show has been published
	p, err := backend.GetProduction(ctx, show.GUID)
	if err != nil || p == nil {
		platform.ReportError(err)
		return nil, err
	}
	if !p.Published {
		return nil, nil // Nope, can't access as it's not public yet
	}

	// list all episodes, excluding future (i.e. unpublished) ones, descending order
	// FIXME filter for other flags, e.g. Block = true
	var er []*podops.Resource
	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(backend.DatastoreResources).Filter("ParentGUID =", show.GUID).Filter("Kind =", podops.ResourceEpisode).Filter("Published <", now).Filter("Published >", 0).Order("-Published").Limit(limit), &er); err != nil {
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

	// track api access for billing etc
	platform.Meter(ctx, "graphql.show", "production", p.GUID)

	return show, nil
}

func (r *queryResolver) Episode(ctx context.Context, guid *string) (*model.Episode, error) {

	data, err := r.EpisodeLoader.Load(ctx, *guid)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}

	// verify that the episode has been published
	episode := data.(*model.Episode)
	published, err := strconv.ParseInt(episode.Published, 10, 64)
	if err != nil {
		platform.ReportError(err)
		return nil, err
	}
	if published == 0 || published > timestamp.Now() {
		return nil, nil // Nope, can't access as it's not public yet
	}

	// track api access for billing etc
	platform.Meter(ctx, "graphql.episode", "episode", *guid)

	return episode, nil
}

func (r *queryResolver) Recent(ctx context.Context, limit int) ([]*model.Show, error) {
	var sh []*podops.Production
	var shows []*model.Show

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(backend.DatastoreProductions).Filter("BuildDate >", 0).Order("-BuildDate").Limit(limit), &sh); err != nil {
		platform.ReportError(err)
		return nil, err
	}

	if len(sh) == 0 {
		shows = make([]*model.Show, 0)
		return shows, nil
	}

	shows = make([]*model.Show, len(sh))
	for i := range sh {
		show, err := r.ShowLoader.Load(ctx, sh[i].Name)
		if err != nil {
			platform.ReportError(err)
			return nil, err
		}
		shows[i] = show.(*model.Show)
	}

	// track api access for billing etc
	platform.Meter(ctx, "graphql.recent", "limit", fmt.Sprintf("%d", limit))

	return shows, nil
}

func (r *queryResolver) Popular(ctx context.Context, limit int) ([]*model.Show, error) {
	return r.Recent(ctx, limit) // FIXME this is just a placeholder, we don't have usage data at the moment to return a real answer
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
