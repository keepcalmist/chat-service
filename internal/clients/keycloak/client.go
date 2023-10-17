package keycloakclient

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

//go:generate options-gen -out-filename=client_options.gen.go -from-struct=Options
type Options struct {
	basePath  string `option:"mandatory" validate:"required,url"`
	realm     string `option:"mandatory" validate:"required"`
	clientID  string `option:"mandatory" validate:"required"`
	secret    string `option:"mandatory" validate:"required"`
	debugMode bool   `option:"default=false"`
}

// Client is a tiny client to the KeyCloak realm operations. UMA configuration:
// http://localhost:3010/realms/Bank/.well-known/uma2-configuration
type Client struct {
	realm    string
	clientID string
	secret   string
	cli      *resty.Client
}

func New(opts Options) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	cli := resty.New()

	cli.SetDebug(opts.debugMode)
	cli.SetBaseURL(opts.basePath)

	return &Client{
		realm:    opts.realm,
		secret:   opts.secret,
		clientID: opts.clientID,
		cli:      cli,
	}, nil
}
