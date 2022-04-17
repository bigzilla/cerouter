package cerouter_test

import (
	"context"
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/injustease/cerouter"
)

func ExampleWithType() {
	r := cerouter.New(cerouter.WithType())
	r.Handle("com.example.ping", func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		fmt.Println(e)
		return nil, nil
	})

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.StartReceiver(context.TODO(), r.Receiver()))
}

func ExampleWithSource() {
	r := cerouter.New(cerouter.WithSource())
	r.Handle("github.com/injustease/cerouter", func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		fmt.Println(e)
		return nil, nil
	})

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.StartReceiver(context.TODO(), r.Receiver()))
}

func ExampleWithSubject() {
	r := cerouter.New(cerouter.WithSubject())
	r.Handle("cerouter", func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		fmt.Println(e)
		return nil, nil
	})

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.StartReceiver(context.TODO(), r.Receiver()))
}

func ExampleWithExtension() {
	r := cerouter.New(cerouter.WithExtension("extkey"))
	r.Handle("extval", func(ctx context.Context, e cloudevents.Event) (*cloudevents.Event, error) {
		fmt.Println(e)
		return nil, nil
	})

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.StartReceiver(context.TODO(), r.Receiver()))
}
