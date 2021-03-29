package backend

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/fupas/commons/pkg/util"
	"github.com/fupas/platform/pkg/platform"

	a "github.com/podops/podops/apiv1"
	"github.com/podops/podops/pkg/backend/models"
)

const (
	// DatastoreProductions collection PRODUCTION
	DatastoreProductions = "PRODUCTIONS"
)

// CreateProduction initializes a new show and all its metadata
func CreateProduction(ctx context.Context, name, title, summary, clientID string) (*models.Production, error) {
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}

	p, err := FindProductionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if p != nil && p.Owner != clientID {
		// do not access someone else's production
		return nil, fmt.Errorf("name '%s' already exists", name)
	}
	if p != nil && p.Owner == clientID {
		// simply return the existing production
		return p, nil
	}

	// create a new production
	id, _ := util.ShortUUID()
	production := strings.ToLower(id)
	now := util.Timestamp()
	location := fmt.Sprintf("%s/show-%s.yaml", production, production)

	p = &models.Production{
		GUID:    production,
		Owner:   clientID,
		Name:    name,
		Title:   title,
		Summary: summary,
		Created: now,
		Updated: now,
	}

	err = UpdateProduction(ctx, p)
	if err != nil {
		return nil, err
	}

	// create a dummy Storage location for this production at production.podops.dev/guid

	show := a.DefaultShow(name, title, summary, production, a.DefaultEndpoint, a.DefaultCDNEndpoint)
	err = WriteResourceContent(ctx, location, true, false, &show)
	if err != nil {
		platform.DataStore().Delete(ctx, productionKey(production))
		return nil, err
	}

	// all done
	return p, nil
}

// GetProduction returns a production based on the GUID
func GetProduction(ctx context.Context, production string) (*models.Production, error) {
	var p models.Production

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
	var p models.Production
	episodes := 0
	show := 0
	assets := 0

	err := platform.DataStore().Get(ctx, productionKey(production), &p)
	if err != nil {
		return err
	}
	rsrc, err := ListResources(ctx, production, a.ResourceALL)
	if err != nil {
		return err
	}
	if rsrc == nil || len(rsrc) == 0 {
		return fmt.Errorf("no resources")
	}

	for _, r := range rsrc {
		if r.Kind == a.ResourceShow {
			show++
		} else if r.Kind == a.ResourceEpisode {
			episodes++
		} else if r.Kind == a.ResourceAsset {
			assets++
		}
	}

	// a podcast needs 1 show and >= 1 episodes to valid
	if show != 1 {
		return fmt.Errorf("missing show")
	}
	if episodes == 0 {
		return fmt.Errorf("missing episodes")
	}
	return nil
}

// UpdateProduction does what the name suggests
func UpdateProduction(ctx context.Context, p *models.Production) error {
	if _, err := platform.DataStore().Put(ctx, productionKey(p.GUID), p); err != nil {
		return err
	}
	return nil
}

// FindProductionByName does a lookup using the productions name instead of its key
func FindProductionByName(ctx context.Context, name string) (*models.Production, error) {
	var p []*models.Production
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreProductions).Filter("Name =", name), &p); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p[0], nil
}

// FindProductionsByOwner returns all productions belonging to the same owner
func FindProductionsByOwner(ctx context.Context, owner string) ([]*models.Production, error) {
	var p []*models.Production
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
