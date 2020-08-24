# lifetime

![Test](https://github.com/TomWright/lifetime/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/TomWright/lifetime)](https://goreportcard.com/report/github.com/TomWright/lifetime)
[![Documentation](https://godoc.org/github.com/TomWright/lifetime?status.svg)](https://godoc.org/github.com/TomWright/lifetime)

Lifetime is a basic package to help you manage the lifetime of an application with multiple routines running at once.

The main benefit of this module is that it allows you to easily manage graceful shutdowns without most of the boilerplate code that goes along with it.

## Installation
```
go get github.com/tomwright/lifetime
```

## Usage

Example usage can be found on godoc.

```
// Create and initialise the lifetime.
lt := lifetime.New(context.Background()).Init()

// Create HTTP server.
mux := http.NewServeMux()
mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte("Hello world"))
})
server := &http.Server{
    Addr:    ":80",
    Handler: mux,
}

// Create HTTP service, giving it the HTTP server.
service := lifetime.NewHTTPService(server)

// Start the service.
lt.Start(service)

go func() {
    // At some point in time call lt.Shutdown.
    // lt.Shutdown would also be executed when a shutdown signal is received.
    time.Sleep(time.Second * 5)
    lt.Shutdown()
}()

// Wait for all services to shutdown
lt.Wait()
```

## Service

A service is a single service within your application that can be started and stopped.

### Graceful shutdown
A graceful shutdown causes all of the `Service.Stop` funcs to be executed causing all services to begin their graceful shutdown.

You can use `lifetime.Wait` to wait for the services to be stopped.

A graceful shutdown will be triggered when:
- A server `Start` func returns an error.
- A `syscall.SIGINT` or `syscall.SIGTERM` signal is received.
- `lifetime.Shutdown` is called.

### Immediate shutdown
An immediate shutdown uses `os.Exit` to immediately stop the application.

This will occur when:
- Multiple `syscall.SIGINT` or `syscall.SIGTERM` signals are received.
- A `syscall.SIGKILL` signal is received.

### Services

Some services are provided for you to use, but you can easily create your own services by implementing the `lifetime.Service` interface.

#### HTTP Server

```
// Create HTTP server.
mux := http.NewServeMux()
mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    writer.Write([]byte("Hello world"))
})
server := &http.Server{
    Addr:    ":80",
    Handler: mux,
}

// Create HTTP service, giving it the HTTP server.
service := lifetime.NewHTTPService(server)

// Start the service.
lt.Start(service)
```

#### GRPC Server

Coming soon...
