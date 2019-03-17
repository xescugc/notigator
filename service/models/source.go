package models

import "github.com/xescugc/notigator/source"

// Source it's the JSON representation of the source.Source
type Source struct {
	Canonical string `json:"canonical"`
	Name      string `json:"name"`
}

// NewSources trasforms src to Source
func NewSources(src []source.Source) []Source {
	rsrc := make([]Source, 0, len(src))

	for _, v := range src {
		rsrc = append(rsrc, Source{Canonical: v.ID, Name: v.Name})
	}

	return rsrc
}
