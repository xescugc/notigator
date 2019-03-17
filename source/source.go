package source

import (
	"fmt"

	"github.com/gosimple/slug"
)

// Canonical it's the type that defines
// the different Sources supported
type Canonical int

//go:generate enumer -type=Canonical -transform=snake

// List of all the supported Sources
const (
	Github Canonical = iota
	Gitlab
	Trello
	Zeplin
)

// Source it's the basic entity representing
// the Source configured
type Source struct {
	ID        string
	Name      string
	Canonical Canonical
}

// BuildID builds the ID of the Source by concatenating
// the 'Canonical' with the 'Name' with a '-'
func (s *Source) BuildID() {
	s.ID = slug.Make(fmt.Sprintf("%s-%s", s.Canonical, s.Name))
}
