package backend

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/datastore"

	ds "github.com/txsvc/platform/v2/pkg/datastore"
	"github.com/txsvc/platform/v2/pkg/timestamp"

	"github.com/podops/podops"
	"github.com/podops/podops/internal/messagedef"
)

// UpdateShow is a helper function to update a show resource
func UpdateShow(ctx context.Context, location string, show *podops.Show) error {
	r, _ := GetResource(ctx, show.GUID())

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != show.Kind {
			return fmt.Errorf(messagedef.MsgResourceKindMismatch, r.Kind, show.Kind)
		}
		r.Name = show.Metadata.Name
		r.Location = location
		r.Title = show.Description.Title
		r.Summary = show.Description.Summary
		r.ImageURI = show.Image.ResolveURI(podops.DefaultStorageEndpoint, show.GUID())
		r.ImageRel = show.Image.Rel
		r.Updated = timestamp.Now()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := timestamp.Now()
	rsrc := podops.Resource{
		Name:       show.Metadata.Name,
		GUID:       show.GUID(),
		Kind:       podops.ResourceShow,
		ParentGUID: show.GUID(),
		Location:   location,
		Title:      show.Description.Title,
		Summary:    show.Description.Summary,
		ImageURI:   show.Image.ResolveURI(podops.DefaultStorageEndpoint, show.GUID()),
		ImageRel:   show.Image.Rel,
		Created:    now,
		Updated:    now,
	}
	return updateResource(ctx, &rsrc)
}

// UpdateEpisode is a helper function to update a episode resource
func UpdateEpisode(ctx context.Context, location string, episode *podops.Episode) error {
	// check if resource with same name already exists for the parent production
	rn, err := FindResource(ctx, episode.Parent(), episode.Metadata.Name)
	if err != nil {
		return err
	}
	r, err := GetResource(ctx, episode.GUID())
	if err != nil {
		return err
	}

	if rn != nil && r != nil {
		if rn.GUID != r.GUID {
			return fmt.Errorf(messagedef.MsgResourceNotFound, fmt.Sprintf("%s/%s", episode.Parent(), episode.Metadata.Name))
		}
	}

	if r != nil {
		// resource already exists, just update the inventory
		if r.Kind != episode.Kind {
			return fmt.Errorf(messagedef.MsgResourceKindMismatch, r.Kind, episode.Kind)
		}
		index, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelEpisode], 10, 64)

		r.Name = episode.Metadata.Name
		r.ParentGUID = episode.Metadata.Labels[podops.LabelParentGUID]
		r.Location = location
		r.Title = episode.Description.Title
		r.Summary = episode.Description.Summary
		r.Published = episode.PublishDateTimestamp()
		r.Index = int(index) // episode number
		r.EnclosureURI = episode.Enclosure.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent())
		r.EnclosureRel = episode.Enclosure.Rel
		r.ImageURI = episode.Image.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent())
		r.ImageRel = episode.Image.Rel
		r.Updated = timestamp.Now()

		return updateResource(ctx, r)
	}

	// create a new inventory entry
	now := timestamp.Now()
	index, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelEpisode], 10, 64)

	rsrc := podops.Resource{
		Name:         episode.Metadata.Name,
		GUID:         episode.GUID(),
		Kind:         podops.ResourceEpisode,
		ParentGUID:   episode.Metadata.Labels[podops.LabelParentGUID],
		Location:     location,
		Title:        episode.Description.Title,
		Summary:      episode.Description.Summary,
		Published:    episode.PublishDateTimestamp(),
		Index:        int(index), // episode number
		EnclosureURI: episode.Enclosure.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent()),
		EnclosureRel: episode.Enclosure.Rel,
		ImageURI:     episode.Image.ResolveURI(podops.DefaultStorageEndpoint, episode.Parent()),
		ImageRel:     episode.Image.Rel,
		Created:      now,
		Updated:      now,
	}
	return updateResource(ctx, &rsrc)
}

func ListPublishedEpisodes(ctx context.Context, production string, published int64, limit int) ([]*podops.Resource, error) {
	var episodes []*podops.Resource

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreResources).Filter("ParentGUID =", production).Filter("Kind =", podops.ResourceEpisode).Filter("Published <", published).Filter("Published >", 0).Order("-Published").Limit(limit), &episodes); err != nil {
		// FIXME filter for other flags, e.g. Block = true
		return nil, err
	}
	return episodes, nil
}

func ListRecentProductions(ctx context.Context, limit int) ([]*podops.Production, error) {
	var shows []*podops.Production

	if _, err := ds.DataStore().GetAll(ctx, datastore.NewQuery(datastoreProductions).Filter("BuildDate >", 0).Order("-BuildDate").Limit(limit), &shows); err != nil {
		return nil, err
	}
	return shows, nil
}

func ListPopularProductions(ctx context.Context, limit int) ([]*podops.Production, error) {
	return ListRecentProductions(ctx, limit) // FIXME this is just a placeholder, we don't have usage data at the moment to return a real answer
}
