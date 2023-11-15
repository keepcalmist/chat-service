package jobsrepo

import "github.com/keepcalmist/chat-service/internal/store"

//go:generate options-gen -out-filename=repo.gen.go -from-struct=Options
type Options struct {
	db *store.Database `option:"mandatory" validate:"required"`
}

type Repo struct {
	Options
}

func New(opts Options) (*Repo, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Repo{Options: opts}, nil
}
