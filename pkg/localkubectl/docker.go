package localkubectl

import (
	"errors"
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

const (
	// LocalkubeLabel is the label that identifies localkube containers
	LocalkubeLabel = "rsprd.com/name=localkube"

	// LocalkubeContainerName is the name of the container that localkube runs in
	LocalkubeContainerName = "localkube"

	// LocalkubeImageName is the image of localkube that is started
	LocalkubeImageName = "redspreadapps/localkube"

	// ContainerDataDir is the path inside the container for etcd data
	ContainerDataDir = "/var/localkube/data"

	// RedspreadName is a Redspread specific identifier for Name
	RedspreadName = "rsprd.com/name"
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

// localkubeCtr returns the localkube container
func (d *Docker) localkubeCtr() (ctrId string, running bool, err error) {
	ctrs, err := d.ListContainers(docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"label": {LocalkubeLabel},
		},
	})

	if err != nil {
		return "", false, fmt.Errorf("Could not list containers: %v", err)
	} else if len(ctrs) < 1 {
		return "", false, ErrNoContainer
	} else if len(ctrs) > 1 {
		return "", false, ErrTooManyLocalkubes
	}

	ctr := ctrs[0]
	return ctr.ID, runningStatus(ctr.Status), nil
}

// createCtr creates the localkube container
func (d *Docker) createCtr(name, imageTag string) (ctrId string, running bool, err error) {
	image := fmt.Sprintf("%s:%s", LocalkubeImageName, imageTag)
	ctrOpts := docker.CreateContainerOptions{
		Name: LocalkubeContainerName,
		Config: &docker.Config{
			Hostname: name,
			Image:    image,
			Env: []string{
				fmt.Sprintf("KUBE_ETCD_DATA_DIRECTORY=%s", ContainerDataDir),
			},
			Labels: map[string]string{
				RedspreadName: name,
			},
			StopSignal: "SIGINT",
		},
	}

	ctr, err := d.CreateContainer(ctrOpts)
	if err != nil {
		if err == docker.ErrNoSuchImage {
			// if image does not exist, pull it
			if pullErr := d.pullImage(imageTag); pullErr != nil {
				return "", false, pullErr
			}
			return d.createCtr(name, imageTag)
		}
		return "", false, fmt.Errorf("Could not create locakube container: %v", err)
	}
	return ctr.ID, ctr.State.Running, nil
}

// pullImage will pull the localkube image on the connected Docker daemon
func (d *Docker) pullImage(imageTag string) error {
	pullOpts := docker.PullImageOptions{
		Repository: LocalkubeImageName,
		Tag:        imageTag,
	}
	err := d.PullImage(pullOpts, docker.AuthConfiguration{})
	if err != nil {
		return fmt.Errorf("Failed to pull localkube image: %v", err)
	}
	return nil
}

// runningStatus returns true if a Docker status string indicates the container is running
func runningStatus(status string) bool {
	return strings.HasPrefix(status, "Up")
}

var (
	// ErrNoContainer is returned when the localkube container hasn't been created yet
	ErrNoContainer = errors.New("Localkube container doesn't exist")

	// ErrTooManyLocalkubes is returned when there are more than one localkube containers on the Docker daemon
	ErrTooManyLocalkubes = errors.New("Multiple localkube containers have been started")
)
