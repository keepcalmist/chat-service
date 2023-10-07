package logger

import (
	"crypto/tls"
	"net/http"

	"github.com/certifi/gocertifi"
	"github.com/getsentry/sentry-go"
)

func NewSentryClient(dsn, env, version string) (*sentry.Client, error) {
	cert, err := gocertifi.CACerts()
	if err != nil {
		return nil, err
	}

	return sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Release:     version,
		Environment: env,
		HTTPTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
			},
		},
		CaCerts: cert,
	})
}
