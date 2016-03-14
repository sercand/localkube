package localkube

import (
	"fmt"
)

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
func (servers Servers) Get(name string) (Server, error) {
	for _, server := range servers {
		if server.Name() == name {
			return server, nil
		}
	}
	return nil, fmt.Errorf("server '%s' does not exist", name)
}

// StartAll starts all services, starting from 0th item and ascending.
func (servers Servers) StartAll() {
	for _, server := range servers {
		fmt.Printf("Starting %s...\n", server.Name())
		server.Start()
	}
}

// StopAll stops all services, starting with the last item.
func (servers Servers) StopAll() {
	for i := len(servers) - 1; i >= 0; i-- {
		server := servers[i]
		fmt.Printf("Stopping %s...\n", server.Name())
		server.Stop()
	}
}

// Start is a helper method to start the Server specified, returns error if server doesn't exist.
func (servers Servers) Start(serverName string) error {
	server, err := servers.Get(serverName)
	if err != nil {
		return err
	}

	server.Start()
	return nil
}

// Stop is a helper method to start the Server specified, returns error if server doesn't exist.
func (servers Servers) Stop(serverName string) error {
	server, err := servers.Get(serverName)
	if err != nil {
		return err
	}

	server.Stop()
	return nil
}

// Status returns a map with the Server name as the key and it's Status as the value.
func (servers Servers) Status() (statuses map[string]Status) {
	for _, server := range servers {
		statuses[server.Name()] = server.Status()
	}
	return statuses
}

// SimpleServer provides a minimal implementation of Server.
type SimpleServer struct {
	ComponentName string
	StartupFn     func()
	ShutdownFn    func()
	StatusFn      func() Status
}

// NoShutdown sets the ShutdownFn to print an error when the server gets shutdown. It returns itself to be chainable.
func (s SimpleServer) NoShutdown() *SimpleServer {
	s.ShutdownFn = func() {
		fmt.Printf("The server '%s' is unstoppable.\n", s.ComponentName)
	}
	return &s
}

// Start calls startup function.
func (s *SimpleServer) Start() {
	s.StartupFn()
}

// Stop calls shutdown function.
func (s *SimpleServer) Stop() {
	if s.ShutdownFn != nil {
		s.ShutdownFn()
	}
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
type Status string

const (
	// Stopped indicates the server is not running.
	Stopped Status = "Stopped"

	// Started indicates the server is running.
	Started = "Started"

	// NotImplemented is returned when Status cannot be determined.
	NotImplemented = "NotImplemented"
)
