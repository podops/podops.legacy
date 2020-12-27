package resources

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/pkg/metadata"
)

const (
	// DatastoreProductions collection PRODUCTION
	DatastoreProductions = "PRODUCTIONS"

	bucketUpload     = "upload.podops.dev" // FIXME make this configurable, otherwise deployment to other locations/projects will not work
	bucketProduction = "production.podops.dev"
	bucketCDN        = "cdn.podops.dev"
)

type (
	// Production holds the shows main data
	Production struct {
		GUID      string `json:"guid"`
		Owner     string `json:"owner"`
		Name      string `json:"name"`
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		Feed      string `json:"feed"`
		NewFeed   string `json:"newFeed"`
		PubDate   int64  `json:"pub_date"`
		BuildDate int64  `json:"build_date"`
		// internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}
)

// CreateProduction initializes a new show and all its metadata
func CreateProduction(ctx context.Context, name, title, summary, clientID string) (*Production, error) {
	if name == "" {
		return nil, fmt.Errorf("production: name must not be empty")
	}

	p, err := FindProductionByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Name '%s' already exists", name))
	}

	id, _ := util.ShortUUID()
	guid := strings.ToLower(id)
	now := util.Timestamp()

	p = &Production{
		GUID:    guid,
		Owner:   clientID,
		Name:    name,
		Title:   title,
		Summary: summary,
		PubDate: now,
		Created: now,
		Updated: now,
	}
	k := productionKey(guid)
	_, err = platform.DataStore().Put(ctx, k, p)
	if err != nil {
		return nil, err
	}

	// create a dummy Storage location for this production at production.podops.dev/guid

	show := metadata.DefaultShow(name, title, summary, guid)
	err = WriteResource(ctx, fmt.Sprintf("%s/show-%s.yaml", guid, guid), true, false, &show)
	if err != nil {
		platform.DataStore().Delete(ctx, k)
		return nil, err
	}

	// all done
	return p, nil
}

// GetProduction returns a production based on the GUID
func GetProduction(ctx context.Context, guid string) (*Production, error) {
	var p Production
	k := productionKey(guid)

	if err := platform.DataStore().Get(ctx, k, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// FindProductionByName does a lookup using the productions name instead of its key
func FindProductionByName(ctx context.Context, name string) (*Production, error) {
	var p []*Production
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreProductions).Filter("Name =", name), &p); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p[0], nil
}

// FindProductionsByOwner returns all productions belonging to the same owner
func FindProductionsByOwner(ctx context.Context, owner string) ([]*Production, error) {
	var p []*Production
	if _, err := platform.DataStore().GetAll(ctx, datastore.NewQuery(DatastoreProductions).Filter("Owner =", owner), &p); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}
	return p, nil
}

func productionKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreProductions, guid, nil)
}
