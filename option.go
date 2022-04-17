package cerouter

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
)

// Option configures the Router.
type Option func(*Router)

// WithType configures the Router to filter by type.
func WithType() Option {
	return func(r *Router) {
		r.filterFn = func(e cloudevents.Event) string {
			return e.Type()
		}
	}
}

// WithSource configures the Router to filter by source.
func WithSource() Option {
	return func(r *Router) {
		r.filterFn = func(e cloudevents.Event) string {
			return e.Source()
		}
	}
}

// WithSubject configures the Router to filter by subject.
func WithSubject() Option {
	return func(r *Router) {
		r.filterFn = func(e cloudevents.Event) string {
			return e.Subject()
		}
	}
}

// WithExtension configures the Router to filter by extension.
func WithExtension(ext string) Option {
	return func(r *Router) {
		r.filterFn = func(e cloudevents.Event) string {
			val, err := types.ToString(e.Extensions()[ext])
			if err != nil {
				return ""
			}
			return val
		}
	}
}
