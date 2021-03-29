package apiv1

import (
	"regexp"

	"github.com/podops/podops/pkg/validator"
)

var (
	nameRegex = regexp.MustCompile(`^[a-z]+[a-z0-9_-]`)
)

// Validate verifies the integrity of struct Show
//
//	APIVersion  string          `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
//	Kind        string          `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: show
//	Metadata    Metadata        `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
//	Description ShowDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
//	Image       Resource        `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'channel.itunes.image'
func (s *Show) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringError(s.APIVersion, Version)
	v.AssertStringError(s.Kind, ResourceShow)

	// Show specific metadata, tracking the scaffolding functions
	v.Validate(&s.Metadata)
	v.AssertContains(s.Metadata.Labels, LabelLanguage, "Metadata")
	v.AssertISO639(s.Metadata.Labels[LabelLanguage])
	v.AssertContains(s.Metadata.Labels, LabelExplicit, "Metadata")
	v.AssertContains(s.Metadata.Labels, LabelType, "Metadata")
	v.AssertContains(s.Metadata.Labels, LabelBlock, "Metadata")
	v.AssertContains(s.Metadata.Labels, LabelComplete, "Metadata")
	v.AssertContains(s.Metadata.Labels, LabelGUID, "Metadata")
	v.Validate(&s.Description)
	v.Validate(&s.Image)

	return v
}

// Validate verifies the integrity of struct Episode
//
//	APIVersion  string             `json:"apiVersion" yaml:"apiVersion" binding:"required"`   // REQUIRED default: v1.0
//	Kind        string             `json:"kind" yaml:"kind" binding:"required"`               // REQUIRED default: episode
//	Metadata    Metadata           `json:"metadata" yaml:"metadata" binding:"required"`       // REQUIRED
//	Description EpisodeDescription `json:"description" yaml:"description" binding:"required"` // REQUIRED
//	Image       Resource           `json:"image" yaml:"image" binding:"required"`             // REQUIRED 'item.itunes.image'
//	Enclosure   Resource           `json:"enclosure" yaml:"enclosure" binding:"required"`     // REQUIRED
func (e *Episode) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringError(e.APIVersion, Version)
	v.AssertStringError(e.Kind, ResourceEpisode)

	// Episode specific metadata, tracking the scaffolding functions
	v.Validate(&e.Metadata)
	v.AssertContains(e.Metadata.Labels, LabelGUID, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelParentGUID, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelDate, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelSeason, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelEpisode, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelExplicit, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelType, "Metadata")
	v.AssertContains(e.Metadata.Labels, LabelBlock, "Metadata")
	v.Validate(&e.Description)
	v.Validate(&e.Image)
	v.Validate(&e.Enclosure)

	return v
}

// Validate verifies the integrity of struct Metadata
//
//	Name   string            `json:"name" yaml:"name" binding:"required"` // REQUIRED <unique name>
//	Labels map[string]string `json:"labels" yaml:"labels,omitempty"`      // REQUIRED
func (m *Metadata) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(m.Name, "Name")
	v.AssertNotEmpty(m.Labels, "Labels")

	return v
}

// Validate verifies the integrity of struct Resource
//
//	URI    string `json:"uri" yaml:"uri" binding:"required"`        // REQUIRED
//	Title  string `json:"title,omitempty" yaml:"title,omitempty"`   // OPTIONAL
//	Anchor string `json:"anchor,omitempty" yaml:"anchor,omitempty"` // OPTIONAL
//	Rel    string `json:"rel,omitempty" yaml:"rel,omitempty"`       // OPTIONAL
//	Type   string `json:"type,omitempty" yaml:"type,omitempty"`     // OPTIONAL
//	Size   int    `json:"size,omitempty" yaml:"size,omitempty"`     // OPTIONAL
func (r *Asset) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(r.URI, "URI")

	return v
}

// Validate verifies the integrity of struct Category
//
//	Name        string   `json:"name" yaml:"name" binding:"required"`      // REQUIRED
//	SubCategory []string `json:"subcategory" yaml:"subcategory,omitempty"` // OPTIONAL
func (c *Category) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(c.Name, "Name")

	return v
}

// Validate verifies the integrity of struct Owner
//
//	Name  string `json:"name" yaml:"name" binding:"required"`   // REQUIRED
//	Email string `json:"email" yaml:"email" binding:"required"` // REQUIRED
func (o *Owner) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(o.Name, "Name")
	v.AssertStringExists(o.Email, "EMail")

	return v
}

// Validate verifies the integrity of struct ShowDescription
//
//	Title     string    `json:"title" yaml:"title" binding:"required"`          // REQUIRED 'channel.title' 'channel.itunes.title'
//	Summary   string    `json:"summary" yaml:"summary" binding:"required"`      // REQUIRED 'channel.description'
//	Link      Resource  `json:"link" yaml:"link"`                               // RECOMMENDED 'channel.link'
//	Category  Category  `json:"category" yaml:"category" binding:"required"`    // REQUIRED channel.category
//	Owner     Owner     `json:"owner" yaml:"owner"`                             // RECOMMENDED 'channel.itunes.owner'
//	Author    string    `json:"author" yaml:"author"`                           // RECOMMENDED 'channel.itunes.author'
//	Copyright string    `json:"copyright,omitempty" yaml:"copyright,omitempty"` // OPTIONAL 'channel.copyright'
//	NewFeed   *Resource `json:"newFeed,omitempty" yaml:"newFeed,omitempty"`     // OPTIONAL channel.itunes.new-feed-url -> move to label
func (d *ShowDescription) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(d.Title, "Title")
	v.AssertStringExists(d.Summary, "Summary")
	v.Validate(&d.Link)
	v.Validate(&d.Category)
	v.Validate(&d.Owner)

	return v
}

// Validate verifies the integrity of struct EpisodeDescription
//
//	Title       string   `json:"title" yaml:"title" binding:"required"`                                 // REQUIRED 'item.title' 'item.itunes.title'
//	Summary     string   `json:"summary" yaml:"summary" binding:"required"`                             // REQUIRED 'item.description'
//	EpisodeText string   `json:"episodeText,omitempty" yaml:"episodeText,omitempty" binding:"required"` // REQUIRED 'item.itunes.summary'
//	Link        Resource `json:"link" yaml:"link"`                                                      // RECOMMENDED 'item.link'
//	Duration    int      `json:"duration" yaml:"duration" binding:"required"`                           // REQUIRED 'item.itunes.duration'
func (d *EpisodeDescription) Validate(v *validator.Validator) *validator.Validator {
	v.AssertStringExists(d.Title, "Title")
	v.AssertStringExists(d.Summary, "Summary")
	v.AssertStringExists(d.EpisodeText, "EpisodeText")
	v.Validate(&d.Link)
	v.AssertNotZero(d.Duration, "Duration")

	return v
}

// ValidResourceName verifies that a name is valid for a resource. The following rules apply:
//
// 'name' must contain only lowercase letters, numbers, dashes (-), underscores (_).
// 'name' must contain 8-44 characters.
// Spaces and dots (.) are not allowed.
func ValidResourceName(name string) bool {
	if len(name) < 8 || len(name) > 45 {
		return false
	}
	return nameRegex.MatchString(name)
}
