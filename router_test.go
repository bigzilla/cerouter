package cerouter_test

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/injustease/cerouter"
	"github.com/injustease/is"
)

func handlePing(is *is.Is) cerouter.Handler {
	return func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		var data string
		is.NoError(e.DataAs(&data))
		is.Equal(data, "ping")

		pongEvent, err := newPongEvent()
		is.NoError(err)

		return pongEvent, nil
	}
}

func handlePong(is *is.Is) cerouter.Handler {
	return func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		var data string
		is.NoError(e.DataAs(&data))
		is.Equal(data, "pong")

		return nil, nil
	}
}

func newPingEvent() (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()

	event.SetID(uuid.New().String())
	event.SetType("com.example.ping")
	event.SetSource("#ping")
	event.SetSubject("handlePing")
	event.SetExtension("myext", "extPing") // string extension
	err := event.SetData(cloudevents.ApplicationJSON, "ping")

	return &event, err
}

func newPongEvent() (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()

	event.SetID(uuid.New().String())
	event.SetType("com.example.pong")
	event.SetSource("#pong")
	event.SetSubject("handlePong")
	event.SetExtension("myext", 1) // non string extension
	err := event.SetData(cloudevents.ApplicationJSON, "pong")

	return &event, err
}

func TestDefaultRouter(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New()
	r.Handle("com.example.ping", handlePing(is))
	r.Handle("com.example.pong", handlePong(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	respEvent, err := r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
	_, err = r.Receiver()(ctx, *respEvent)
	is.NoError(err)
}

func TestWithType(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New(cerouter.WithType())
	r.Handle("com.example.ping", handlePing(is))
	r.Handle("com.example.pong", handlePong(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	respEvent, err := r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
	_, err = r.Receiver()(ctx, *respEvent)
	is.NoError(err)
}

func TestWithSource(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New(cerouter.WithSource())
	r.Handle("#ping", handlePing(is))
	r.Handle("#pong", handlePong(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	respEvent, err := r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
	_, err = r.Receiver()(ctx, *respEvent)
	is.NoError(err)
}

func TestWithSubject(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New(cerouter.WithSubject())
	r.Handle("handlePing", handlePing(is))
	r.Handle("handlePong", handlePong(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	respEvent, err := r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
	_, err = r.Receiver()(ctx, *respEvent)
	is.NoError(err)
}

func TestWithExtension(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New(cerouter.WithExtension("myext"))
	r.Handle("extPing", handlePing(is))
	r.Handle("extPong", handlePong(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	respEvent, err := r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
	_, err = r.Receiver()(ctx, *respEvent)
	is.NoError(err)
}

func TestUnmatchedFilter(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	r := cerouter.New()
	r.Handle("com.example.unknown.ping", handlePing(is))

	pingEvent, err := newPingEvent()
	is.NoError(err)

	_, err = r.Receiver()(ctx, *pingEvent)
	is.NoError(err)
}

func TestHandleValidation(t *testing.T) {
	is := is.New(t)

	var r cerouter.Router
	is.Panic(func() {
		r.Handle("", handlePing(is))
	})
	is.Panic(func() {
		r.Handle("com.example.ping", nil)
	})
	r.Handle("com.example.ping", handlePing(is))
	is.Panic(func() {
		r.Handle("com.example.ping", handlePing(is))
	})
}

func TestReceiverValidation(t *testing.T) {
	is := is.New(t)
	ctx := context.TODO()

	var r cerouter.Router

	_, err := r.Receiver()(ctx, cloudevents.Event{})
	is.Error(err)
}
