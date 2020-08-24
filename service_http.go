package lifetime

import "net/http"

// NewHTTPService returns a service that will run listen and serve the given
// HTTP server.
func NewHTTPService(server *http.Server) Service {
	return &httpService{
		server: server,
	}
}

// httpService is an implementation of Service that will listen and serve the given
// HTTP server.
type httpService struct {
	server *http.Server
}

// Start will start the service.
// This is a blocking call and should block for the lifetime of the service.
// Returns an error which is treated as fatal.
func (service *httpService) Start() error {
	err := service.server.ListenAndServe()
	if err == nil {
		return nil
	}
	// ErrServerClosed is returned when we call service.Close() from Service.Stop
	// so we shouldn't treat it as a breaking error.
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Stop will stop the service.
// Stop is not called if Start returned an error.
func (service *httpService) Stop() {
	_ = service.server.Close()
}
