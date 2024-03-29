// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Category struct {
	Name        string  `json:"name"`
	Subcategory *string `json:"subcategory"`
}

type Enclosure struct {
	Link string `json:"link"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Episode struct {
	GUID        string              `json:"guid"`
	Name        string              `json:"name"`
	Created     string              `json:"created"`
	Published   string              `json:"published"`
	Labels      *Labels             `json:"labels"`
	Description *EpisodeDescription `json:"description"`
	Image       string              `json:"image"`
	Enclosure   *Enclosure          `json:"enclosure"`
	Production  *Production         `json:"production"`
}

type EpisodeDescription struct {
	Title       string  `json:"title"`
	Summary     string  `json:"summary"`
	Description *string `json:"description"`
	Link        string  `json:"link"`
	Duration    int     `json:"duration"`
}

type Labels struct {
	Block    string `json:"block"`
	Explicit string `json:"explicit"`
	Type     string `json:"type"`
	Complete string `json:"complete"`
	Language string `json:"language"`
	Episode  int    `json:"episode"`
	Season   int    `json:"season"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Production struct {
	GUID  string `json:"guid"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

type Show struct {
	GUID        string           `json:"guid"`
	Name        string           `json:"name"`
	Created     string           `json:"created"`
	Build       string           `json:"build"`
	Labels      *Labels          `json:"labels"`
	Description *ShowDescription `json:"description"`
	Image       string           `json:"image"`
	Episodes    []*Episode       `json:"episodes"`
}

type ShowDescription struct {
	Title     string      `json:"title"`
	Summary   string      `json:"summary"`
	Link      string      `json:"link"`
	Category  []*Category `json:"category"`
	Author    string      `json:"author"`
	Copyright string      `json:"copyright"`
	Owner     *Owner      `json:"owner"`
}
