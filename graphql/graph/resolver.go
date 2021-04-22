package graph

import (
	"context"
	"fmt"
	"strconv"

	"github.com/podops/podops"
	"github.com/podops/podops/backend"
	"github.com/podops/podops/graphql/graph/model"
	"github.com/podops/podops/internal/loader"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver holds the loaders
type Resolver struct {
	ShowLoader    *loader.Loader
	EpisodeLoader *loader.Loader
}

// LoadShow loads a show
func LoadShow(ctx context.Context, key string) (interface{}, error) {
	p, err := backend.FindProductionByName(ctx, key)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("show '%s' not found", key)
	}

	r, err := backend.GetResource(ctx, p.GUID)
	if err != nil {
		return nil, err
	}

	s, err := backend.GetResourceContent(ctx, p.GUID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("show '%s': no content", key)
	}
	show := s.(*podops.Show)

	category := make([]*model.Category, 1)
	category[0] = &model.Category{
		Name:        show.Description.Category.Name,
		Subcategory: &show.Description.Category.SubCategory[0],
	}

	labels := &model.Labels{
		Block:    show.Metadata.Labels[podops.LabelBlock],
		Explicit: show.Metadata.Labels[podops.LabelExplicit],
		Type:     show.Metadata.Labels[podops.LabelType],
		Complete: show.Metadata.Labels[podops.LabelComplete],
		Language: show.Metadata.Labels[podops.LabelLanguage],
		//Episode: NOT USED
		//Season: NOT USED
	}

	return &model.Show{
		GUID:    p.GUID,
		Name:    p.Name,
		Created: strconv.FormatInt(p.Created, 10),
		Build:   strconv.FormatInt(p.BuildDate, 10),
		Labels:  labels,
		Description: &model.ShowDescription{
			Title:     show.Description.Title,
			Summary:   show.Description.Summary,
			Link:      show.Description.Link.URI,
			Category:  category,
			Author:    show.Description.Author,
			Copyright: show.Description.Copyright,
			Owner: &model.Owner{
				Name:  show.Description.Owner.Name,
				Email: show.Description.Owner.Email,
			},
		},
		Image: r.ImageURI,
		// Episodes are loaded by hthe schema.resolver implementation in order make use of the dataloader
	}, nil
}

// LoadEpisode loads an episode
func LoadEpisode(ctx context.Context, key string) (interface{}, error) {
	r, err := backend.GetResource(ctx, key)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, fmt.Errorf("episode '%s' not found", key)
	}

	meta, err := backend.GetMetadataForResource(ctx, key)
	if err != nil {
		return nil, err
	}

	e, err := backend.GetResourceContent(ctx, r.GUID)
	if err != nil {
		return nil, err
	}
	episode := e.(*podops.Episode)

	p, err := backend.GetProduction(ctx, r.ParentGUID)
	if err != nil {
		return nil, err
	}

	n, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelEpisode], 10, 64)
	season, _ := strconv.ParseInt(episode.Metadata.Labels[podops.LabelSeason], 10, 64)
	labels := &model.Labels{
		Block:    episode.Metadata.Labels[podops.LabelBlock],
		Explicit: episode.Metadata.Labels[podops.LabelExplicit],
		Type:     episode.Metadata.Labels[podops.LabelType],
		Complete: episode.Metadata.Labels[podops.LabelComplete],
		//Language: NOT USED
		Episode: int(n),
		Season:  int(season),
	}

	return &model.Episode{
		GUID:      episode.GUID(),
		Name:      episode.Metadata.Name,
		Created:   strconv.FormatInt(r.Created, 10),
		Published: strconv.FormatInt(r.Published, 10),
		Labels:    labels,
		Description: &model.EpisodeDescription{
			Title:       episode.Description.Title,
			Summary:     episode.Description.Summary,
			Description: &episode.Description.EpisodeText,
			Link:        episode.Description.Link.URI,
			Duration:    int(meta.Duration), // episode.Description.Duration,
		},
		Image: r.ImageURI,
		Enclosure: &model.Enclosure{
			Link: r.EnclosureURI, //episode.Enclosure.URI,
			Type: r.EnclosureRel, // episode.Enclosure.Type,
			Size: int(meta.Size), // episode.Enclosure.Size,
		},
		Production: &model.Production{
			GUID:  p.GUID,
			Name:  p.Name,
			Title: p.Title,
		},
	}, nil
}

// CreateResolver returns a resolver for loading shows and episodes
func CreateResolver() *Resolver {
	return &Resolver{
		ShowLoader:    loader.New(LoadShow, loader.DefaultTTL),
		EpisodeLoader: loader.New(LoadEpisode, loader.DefaultTTL),
	}
}
