package source

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Canonical int

type Source struct {
	ID        string
	Name      string
	Canonical Canonical
}

func (s *Source) BuildID() {
	s.ID = slug.Make(fmt.Sprintf("%s-%s", s.Canonical, s.Name))
}

//go:generate enumer -type=Canonical -transform=snake

const (
	Github Canonical = iota
	Gitlab
	Trello
	Zeplin
)
