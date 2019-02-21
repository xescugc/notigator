package models

import "github.com/xescugc/notigator/source"

type Source struct {
	Canonical string `json:"canonical"`
	Name      string `json:"name"`
}

func NewSources(src []source.Source) []Source {
	rsrc := make([]Source, 0, len(src))

	for _, v := range src {
		rsrc = append(rsrc, Source{Canonical: v.Canonical.String(), Name: v.Name})
	}

	return rsrc
}
