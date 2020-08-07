package lifetime

// Service defines a single service in an application.
type Service interface {
	// Start will start the service.
	// This is a blocking call and should block for the lifetime of the service.
	// Returns an error which is treated as fatal.
	Start() error
	// Stop will stop the service.
	// Stop is not called if Start returned an error.
	Stop()
}
