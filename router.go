package cerouter

import (
	"context"
	"sync"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Handler is the function signature for CloudEvent receiver.
type Handler func(context.Context, cloudevents.Event) (*cloudevents.Event, error)

// Router is a CloudEvent router which can dispatch CloudEvent by its context attribute.
// Zero value Router is a ready to use router. Default filter is filter by type.
type Router struct {
	mu       sync.RWMutex
	m        map[string]Handler
	filterFn func(cloudevents.Event) string
}

// New returns an initialized Router and accept optional Option.
// New without Option is the same as using zero value Router.
// Configure filter by other context attribute using Option, for example.
//
// 	r := cerouter.New(cerouter.WithSource())
func New(opts ...Option) *Router {
	r := &Router{
		m: make(map[string]Handler),
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.filterFn == nil {
		WithType()(r)
	}

	return r
}

// Handle registers the handler for the given filter.
func (r *Router) Handle(filter string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if filter == "" {
		panic("cerouter: invalid filter")
	}
	if handler == nil {
		panic("cerouter: nil handler")
	}

	if r.m == nil {
		r.m = make(map[string]Handler)
	}

	if _, ok := r.m[filter]; ok {
		panic("cerouter: multiple registrations for " + filter)
	}

	r.m[filter] = handler
}

// Receiver returns a Handler to be used by CloudEvent receiver.
func (r *Router) Receiver() Handler {
	if r.m == nil {
		r.m = make(map[string]Handler)
	}
	if r.filterFn == nil {
		WithType()(r)
	}

	return func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		if err := e.Validate(); err != nil {
			return nil, err
		}

		if fn, ok := r.m[r.filterFn(e)]; ok {
			return fn(ctx, e)
		}
		return nil, nil
	}
}
