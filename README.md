# cerouter [![Go Reference](https://pkg.go.dev/badge/github.com/injustease/cerouter.svg)](https://pkg.go.dev/github.com/injustease/cerouter) [![main](https://github.com/injustease/cerouter/workflows/main/badge.svg)](https://github.com/injustease/cerouter/actions/workflows/main.yaml)

Package cerouter provides a router for [CloudEvents](https://github.com/cloudevents/sdk-go).

## Get started

Install cerouter.

```bash
go get github.com/injustease/cerouter
```

Import the module into your code.

```go
import "github.com/injustease/cerouter"
```

While the [Go SDK for CloudEvents](https://github.com/cloudevents/sdk-go) provides multiple signature for the [receiver](https://github.com/cloudevents/sdk-go/blob/main/v2/client/client.go#L33-L50), this package only support one of them.

```go
type Handler func(context.Context, cloudevents.Event) (*cloudevents.Event, error)
```

Default router.

```go
func main() {
    r := cerouter.New() // default router will filter event by type
    r.Handle("com.example.ping", handlePing)
    r.Handle("com.example.pong", handlePong)

    c, err := cloudevents.NewClientHTTP()
    if err != nil {
        log.Fatalf("failed to create client, %v", err)
    }

    log.Fatal(c.StartReceiver(context.TODO(), r.Receiver()))
}

func handlePing(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
    // handle ping
}

func handlePong(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
    // handle pong
}
```

You can specify which context attribute to filter.

```go
func main() {
    typeRouter := cerouter.New(cerouter.WithType()) // filter event by type, same as default one
    typeRouter.Handle("com.example.ping", handlePing)

    sourceRouter := cerouter.New(cerouter.WithSource()) // filter event by source
    sourceROuter.Handle("github.com/injustease/cerouter")

    subjectRouter := cerouter.New(cerouter.WithSubject()) // filter event by subject
    subjectRouter.Handle("cerouter", handlePing)

    extensionRouter := cerouter.New(cerouter.WithExtension("extkey")) // filter event by extension
    extensionRouter.Handle("extval", handlePing)
}
```

## Example

Run the following code, it will starting HTTP server with default port `:8080` and default path `/`.

```go
package main

import (
    "context"
    "fmt"
    "log"

    cloudevents "github.com/cloudevents/sdk-go/v2"
    "github.com/injustease/cerouter"
)

func main() {
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
```

Try to send the event.

```bash
curl -v "localhost:8080" \
-X POST \
-H "Ce-Id: unique-id" \
-H "Ce-Specversion: 1.0" \
-H "Ce-Type: com.example.ping" \
-H "Ce-Source: github.com/injustease/cerouter" \
-H "Ce-Subject: cerouter" \
-H "Ce-Extkey: extval" \
-H "Content-Type: application/json" \
-d '{"msg":"Hello from cerouter!"}'
```
