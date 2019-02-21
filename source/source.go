package source

type Canonical int

type Source struct {
	Canonical Canonical
	Name      string
}

//go:generate enumer -type=Canonical -transform=snake

const (
	Github Canonical = iota
	Gitlab
	Trello
)

var (
	Sources = []Source{
		Source{
			Canonical: Github,
			Name:      "GitHub",
		},
		Source{
			Canonical: Gitlab,
			Name:      "GitLab",
		},
		Source{
			Canonical: Trello,
			Name:      "Trello",
		},
	}
)
