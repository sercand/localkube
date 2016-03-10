package localkube

// Server represents a component that Kubernetes depends on. It allows for the management of
// the lifecycle of the component.
type Server interface {
	// Start immediately starts the component.
	Start()
	// Stop begins the process of stopping the components.
	Stop()

	// Name returns a unique identifier for the component.
	Name() string

	// Status provides the state of the server.
	Status() Status
}

// Servers allows operations to be performed on many servers at once.
// Uses slice to preserve ordering.
type Servers []Server

// Get returns a server matching name, returns nil if server doesn't exit.
func (servers Servers) Get(name string) Server {
	for _, server := range servers {
		if server.Name() == name {
			return server
		}
	}
	return nil
}

// StartAll starts all services, starting from 0th item and ascending.
func (servers Servers) StartAll() {
	for _, server := range servers {
		server.Start()
	}
}

// StopAll stops all services, starting with the last item and descending.
func (servers Servers) StopAll() {
	for _, server := range servers {
		server.Stop()
	}
}

// Status returns a map with the Server name as the key and it's Status as the value.
func (servers Servers) Status() (statuses map[string]Status) {
	for _, server := range servers {
		statuses[server.Name()] = server.Status()
	}
	return statuses
}

// SimpleServer provides and easy implementation of Server.
type SimpleServer struct {
	ComponentName string
	StartupFn     func()
	ShutdownFn    func()
	StatusFn      func() Status
}

// Start calls startup function.
func (s *SimpleServer) Start() {
	s.StartupFn()
}

// Stop calls shutdown function.
func (s *SimpleServer) Stop() {
	s.ShutdownFn()
}

// Name returns the name of the service.
func (s SimpleServer) Name() string {
	return s.ComponentName
}

// Status calls the status function and returns the the Server's status.
func (s *SimpleServer) Status() Status {
	return s.StatusFn()
}

// Status indicates the condition of a Server.
type Status uint

const (
	// Stopped indicates the server is not running.
	Stopped Status = iota

	// Started indicates the server is running.
	Started

	// NotImplemented is returned when Status cannot be determined.
	NotImplemented
)
