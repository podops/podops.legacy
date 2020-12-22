package production

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/podops/podops/internal/errors"
)

const (
	// DatastoreProductions collection PRODUCTION
	DatastoreProductions = "PRODUCTIONS"
)

type (
	// Production holds the shows main data
	Production struct {
		GUID      string `json:"guid"`
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
func CreateProduction(ctx context.Context, name, title, summary string) (*Production, error) {

	p, err := FindProductionByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if p != nil {
		return nil, errors.New(fmt.Sprintf("Show with name '%s' already exists", name), http.StatusConflict)
	}

	guid, _ := util.ShortUUID()
	now := util.Timestamp()

	p = &Production{
		GUID:    guid,
		Name:    name,
		Title:   title,
		Summary: summary,
		PubDate: now,
		Created: now,
		Updated: now,
	}
	k := productionKey(guid)
	_, err = platform.DataStore().Put(ctx, k, p)

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

func productionKey(guid string) *datastore.Key {
	return datastore.NameKey(DatastoreProductions, guid, nil)
}
