// Package xrequesttrace a traefik plugin.
package xrequesttrace

import (
	"context"
	"net/http"
	"regexp"
)

// Config the plugin configuration.
type Config struct {
	Dummy bool `json:"dummy,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Dummy: false,
	}
}

// XRequestTrace a XRequestTrace plugin.
type XRequestTrace struct {
	next http.Handler
	name string
}

// New created a new XRequestTrace plugin.
func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &XRequestTrace{
		next: next,
		name: name,
	}, nil
}

var traceRegex, _ = regexp.Compile(`^\w{2}-(\w{32})-\w{16}-\w{2}$`)

func (a *XRequestTrace) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if xids, ok := req.Header["X-Request-ID"]; ok && len(xids) == 0 {
		if traceparents, ok := req.Header["traceparent"]; ok && len(traceparents) > 0 {
			traceparent := traceparents[0]
			if match := traceRegex.MatchString(traceparent); match {
				traceid := traceRegex.FindStringSubmatch(traceparent)[1]
				req.Header.Set("X-Request-ID", traceid)
			}
		}
	}

	a.next.ServeHTTP(rw, req)
}
