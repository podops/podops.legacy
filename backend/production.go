package backend

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"
	"github.com/podops/podops"
	"github.com/podops/podops/internal/errordef"
)

const (
	// DatastoreProductions collection PRODUCTION
	DatastoreProductions = "PRODUCTIONS"
)

// CreateProduction initializes a new show and all its metadata
func CreateProduction(ctx context.Context, name, title, summary, clientID string) (*podops.Production, error) {
	if name == "" {
		return nil, errordef.ErrInvalidParameters
	}

	p, err := FindProductionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if p != nil {
		if p.Owner != clientID {
			// do not access someone else's production
			return nil, fmt.Errorf(errordef.MsgResourceAlreadyExists, name)
		}
		return p, nil
	}

	// create a new production
	id, _ := util.ShortUUID()
	production := strings.ToLower(id)
	now := util.Timestamp()

	prod := podops.Production{
		GUID:    production,
		Owner:   clientID,
		Name:    name,
		Title:   title,
		Summary: summary,
		Created: now,
		Updated: now,
	}

	err = UpdateProduction(ctx, &prod)
	if err != nil {
		return nil, err
	}

	return &prod, nil
}

// GetProduction returns a production based on the GUID
func GetProduction(ctx context.Context, production string) (*podops.Production, error) {
	var p podops.Production

	if err := platform.DataStore().Get(ctx, productionKey(production), &p); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &p, nil
}

// ValidateProduction checks the integrity of a production and fixes issues if possible
func ValidateProduction(ctx context.Context, production string) error {
	var p podops.Production
	episodes := 0
	show := 0
	assets := 0

	err := platform.DataStore().Get(ctx, productionKey(production), &p)
	if err != nil {
		return err
	}
	rsrc, err := ListResources(ctx, production, podops.ResourceALL)
	if err != nil {
		return err
	}
	if len(rsrc) == 0 {
		return errordef.ErrNoSuchResource
	}

	for _, r := range rsrc {
		if r.Kind == podops.ResourceShow {
			show++
		} else if r.Kind == podops.ResourceEpisode {
			episodes++
		} else if r.Kind == podops.ResourceAsset {
			assets++
		}
	}

	// a podcast needs 1 show and >= 1 episodes to valid
	if show != 1 {
		return errordef.ErrNoSuchProduction
	}
	if episodes == 0 {
		return errordef.ErrNoSuchEpisode
	}
	return nil
}

// UpdateProduction does what the name suggests
func UpdateProduction(ctx context.Context, p *podops.Production) error {
	if _, err := platform.DataStore().Put(ctx, productionKey(p.GUID), p); err != nil {
		return err
	}
	return nil
}

// FindProductionByName does a lookup using the productions name instead of its key
func FindProductionByName(ctx context.Context, name string) (*podops.Production, error) {
	var p []*podops.Production
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreProductions).Filter("Name =", name), &p); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p[0], nil
}

// FindProductionsByOwner returns all productions belonging to the same owner
func FindProductionsByOwner(ctx context.Context, owner string) ([]*podops.Production, error) {
	var p []*podops.Production
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreProductions).Filter("Owner =", owner), &p); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p, nil
}

func productionKey(production string) *datastore.Key {
	return datastore.NameKey(DatastoreProductions, production, nil)
}
