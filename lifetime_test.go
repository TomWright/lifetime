package lifetime_test

import (
	"context"
	"fmt"
	"github.com/tomwright/lifetime"
	"time"
)

type testService struct {
	name             string
	stop             bool
	startupDuration  time.Duration
	shutdownDuration time.Duration
}

func (s *testService) Start() error {
	time.Sleep(s.startupDuration)
	fmt.Printf("%s: Started\n", s.name)
	for s.stop == false {
		time.Sleep(time.Millisecond * 10)
	}
	return nil
}

func (s *testService) Stop() {
	time.Sleep(s.shutdownDuration)
	fmt.Printf("%s: Stopped\n", s.name)
	s.stop = true
}

// ExampleLifetime shows a basic example of how you can use Lifetime.
func ExampleLifetime() {
	// Create a lifetime and initialises it.
	lt := lifetime.New(context.Background()).
		Init()

	fmt.Printf("Starting services\n")

	// Service A takes 100ms to start up and 800ms to shutdown.
	serviceA := &testService{
		name:             "a",
		startupDuration:  time.Millisecond * 100,
		shutdownDuration: time.Millisecond * 800,
	}
	// Service B takes 800ms to start up and 100ms to shutdown.
	serviceB := &testService{
		name:             "b",
		startupDuration:  time.Millisecond * 800,
		shutdownDuration: time.Millisecond * 100,
	}

	// Start both services.
	lt.Start(serviceA)
	lt.Start(serviceB)

	// Wait some time and trigger an application shutdown.
	go func() {
		<-time.After(time.Millisecond * 1500)
		fmt.Printf("Shutting down\n")
		lt.Shutdown()
	}()

	// Wait for all services to stop.
	lt.Wait()

	fmt.Printf("Shutdown\n")

	// Output:
	// Starting services
	// a: Started
	// b: Started
	// Shutting down
	// b: Stopped
	// a: Stopped
	// Shutdown
}
