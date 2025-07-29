// Package xrequesttrace a traefik plugin.
package xrequesttrace

import (
	"context"
	"log"
	"net/http"
	"os"
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
	name   string
	next   http.Handler
	logger *log.Logger
}

// New created a new XRequestTrace plugin.
func New(_ context.Context, next http.Handler, _ *Config, name string) (http.Handler, error) {
	return &XRequestTrace{
		name:   name,
		next:   next,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}, nil
}

var traceRegex, _ = regexp.Compile(`^\w{2}-(\w{32})-\w{16}-\w{2}$`)

func (a *XRequestTrace) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	a.logger.Printf("headers received: %#v", req.Header)
	a.logger.Printf("request context: %#v", ctx)
	if xids := req.Header["X-Request-ID"]; len(xids) == 0 {
		a.logger.Printf("X-Request-ID header not found")
		if traceparents := req.Header["traceparent"]; len(traceparents) > 0 {
			traceparent := traceparents[0]
			a.logger.Printf("traceparent header found: %#v", traceparent)
			if match := traceRegex.MatchString(traceparent); match {
				traceid := traceRegex.FindStringSubmatch(traceparent)[1]
				a.logger.Printf("traceid extracted: %#v", traceid)
				req.Header["X-Request-ID"] = []string{traceid}
			} else {
				a.logger.Printf("traceparent header does not match expected format: %#v", traceparent)
			}
		} else {
			a.logger.Printf("traceparent header not found %#v", traceparents)
		}
	} else {
		a.logger.Printf("X-Request-ID header already exists: %#v", xids)
	}

	a.next.ServeHTTP(rw, req)
}
