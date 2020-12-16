package metadata

import (
	"time"

	"github.com/podops/podops/pkg/feed"
)

type (
	Metadata struct {
		Name   string
		Labels map[string]string `yaml:"labels,omitempty"`
	}

	Link struct {
		URI   string
		Title string `yaml:"title,omitempty"`
		Rel   string `yaml:"rel,omitempty"`
		Type  string `yaml:"type,omitempty"`
	}

	Image struct {
		URI    string
		Title  string `yaml:"title,omitempty"`
		Anchor string `yaml:"anchor,omitempty"`
	}

	Contributor struct {
		Name string
		URI  string `yaml:"uri,omitempty"`
	}

	Owner struct {
		Name  string
		Email string
	}

	Description struct {
		Title     string
		Link      string
		SubTitle  string `yaml:"subtitle,omitempty"`
		Summary   string `yaml:"summary,omitempty"`
		Copyright string `yaml:"copyright,omitempty"`
		Author    string
		Owner     Owner
	}

	Content struct {
		Category string
		Type     string
		Language string `yaml:"language,omitempty"`
		Explicit string `yaml:"explicit,omitempty"`
		Block    string `yaml:"block,omitempty"`
	}

	// Show holds all metadata related to a podcast/show
	Show struct {
		APIVersion   string `yaml:"apiVersion"`
		Kind         string
		Metadata     Metadata `yaml:"metadata,omitempty"`
		Description  Description
		Content      Content
		Image        Image
		Links        []Link        `yaml:"links,omitempty"`
		Contributors []Contributor `yaml:"contributors,omitempty"`
	}
)

func (l *Link) ToAtomLink() *feed.AtomLink {
	return &feed.AtomLink{
		HREF: l.URI,
		Rel:  l.Rel,
		Type: l.Type,
	}
}
func (show *Show) ToPodcast() (*feed.Podcast, error) {
	now := time.Now()
	// basics
	p := feed.New(show.Description.Title, show.Description.Link, show.Description.Summary, &now, &now)
	// details
	p.AddSubTitle(show.Description.SubTitle)
	p.AddSummary(show.Description.Summary)
	p.AddAuthor(show.Description.Owner.Name, show.Description.Owner.Email)
	p.AddCategory(show.Content.Category, nil)
	p.AddImage(show.Image.URI)
	/*
		// add atom links
		for _, link := range show.Links {
			p.AddAtomLink(link.URI)
		}
	*/
	return &p, nil
}
