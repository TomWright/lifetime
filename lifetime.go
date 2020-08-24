package lifetime

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	// ErrShutdownSignalReceived is used when a shutdown signal is received.
	// It will cause a graceful shutdown.
	ErrShutdownSignalReceived = errors.New("shutdown signal received")

	// ErrImmediateShutdownSignalReceived is used when a shutdown signal is received for the second time.
	// It will cause an immediate shutdown.
	ErrImmediateShutdownSignalReceived = errors.New("immediate shutdown signal received")
)

// New returns a new Lifetime instance that can be used to control
// the lifetime of an application.
func New(ctx context.Context) *Lifetime {
	ctx, cancel := context.WithCancel(ctx)
	return &Lifetime{
		ctx:        ctx,
		cancelFunc: cancel,
		serviceWg:  &sync.WaitGroup{},
		errCh:      make(chan error),
	}
}

// Lifetime contains some basic functionality you can use to control the lifetime of an application.
type Lifetime struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	serviceWg  *sync.WaitGroup
	errCh      chan error
}

// Init starts up the required routines for the lifetime instance to work as expected.
func (lifetime *Lifetime) Init() *Lifetime {
	lifetime.handleErrors()
	lifetime.handleShutdownSignals()
	return lifetime
}

// Context returns a context that should be used throughout the runtime of the application.
// When a shutdown of the application is triggered this context will be closed.
func (lifetime *Lifetime) Context() context.Context {
	return lifetime.ctx
}

// Done returns a channel that's closed when all work should be stopped.
// This is really just the Done channel of the lifetime context.
func (lifetime *Lifetime) Done() <-chan struct{} {
	return lifetime.ctx.Done()
}

// Shutdown triggers a graceful shutdown of the application.
func (lifetime *Lifetime) Shutdown() {
	lifetime.cancelFunc()
}

// Wait will block until all services registered with the Lifetime have finished execution.
func (lifetime *Lifetime) Wait() {
	lifetime.serviceWg.Wait()
}

// Start will start the given service.
// It also ensures that the service wait group is updated as expected.
func (lifetime *Lifetime) Start(svc Service) {
	lifetime.serviceWg.Add(1)
	go lifetime.start(svc)
}

// start executes a service in a go routine.
// It ensures that the service wait group is updated, and that the service Stop func is
// executed when an application shutdown is triggered.
func (lifetime *Lifetime) start(svc Service) {
	defer lifetime.serviceWg.Done()

	startErrs := make(chan error)
	startWg := &sync.WaitGroup{}

	startWg.Add(1)
	go func() {
		defer startWg.Done()
		err := svc.Start()
		if err != nil {
			startErrs <- err
		}
	}()

	select {
	case startErr := <-startErrs:
		// Something went wrong during start-up.
		// Report the error.
		lifetime.errCh <- startErr
	case <-lifetime.ctx.Done():
		// The application wants us to shutdown.
		// Stop the service and wait for the start func to finish.
		svc.Stop()
		startWg.Wait()
	}
}

// handleShutdownSignals runs a go routine that listens for shutdown signals from the os
// and sends an ErrShutdownSignalReceived to the error chan when the application is told to shutdown.
func (lifetime *Lifetime) handleShutdownSignals() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		count := 0
		for {
			sig := <-signals
			count++
			if count > 1 || sig == syscall.SIGKILL {
				lifetime.errCh <- ErrImmediateShutdownSignalReceived
				continue
			}
			lifetime.errCh <- ErrShutdownSignalReceived
		}
	}()
}

// handleErrors starts a go routine that listens on the error channel and logs errors.
func (lifetime *Lifetime) handleErrors() {
	go func() {
		for {
			err, ok := <-lifetime.errCh
			if !ok {
				lifetime.cancelFunc()
				return
			}

			if err == ErrImmediateShutdownSignalReceived {
				os.Exit(1)
			}

			log.Printf("lifetime error received: %s", err.Error())

			lifetime.Shutdown()
		}
	}()
}
