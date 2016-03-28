package localkubectl

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

const (
	// LocalkubeLabel is the label that identifies localkube containers
	LocalkubeLabel = "rsprd.com/name=localkube"
)

// Docker provides a wrapper around the Docker client for easy control of localkube.
type Docker struct {
	*docker.Client
}

// NewDocker returns a localkube Docker client from a created *docker.Client
func NewDocker(client *docker.Client) (*Docker, error) {
	_, err := client.Version()
	if err != nil {
		return nil, fmt.Errorf("Unable to establish connection with Docker daemon: %v", err)
	}

	return &Docker{
		Client: client,
	}, nil
}

// NewDockerFromEnv creates a new Docker client using environment clues.
func NewDockerFromEnv() (*Docker, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Could not create Docker client: %v", err)
	}

	return NewDocker(client)
}

// localkubeCtrs lists the containers associated with localkube. If running is true, only running containers will be listed.
func (d *Docker) localkubeCtrs(runningOnly bool) ([]docker.APIContainers, error) {
	ctrs, err := d.ListContainers(docker.ListContainersOptions{
		All: !runningOnly,
		Filters: map[string][]string{
			"label": {LocalkubeLabel},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("Could not list containers: %v", err)
	}
	return ctrs, nil
}
