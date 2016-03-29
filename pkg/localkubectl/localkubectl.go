// Package localkubectl allows the lifecycle of the localkube container to be controlled
package localkubectl

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

var (
	// DefaultHostDataDir is the directory on the host which is mounted for localkube data if no directory is specified
	DefaultHostDataDir = "~/.localkube/data"

	// LocalkubeDefaultTag is the tag to use for the localkube image
	LocalkubeDefaultTag = "latest"

	// LocalkubeClusterName is the name the cluster configuration is stored under
	LocalkubeClusterName = "localkube"

	// LocalkubeContext is the name of the context used by localkube
	LocalkubeContext = "localkube"
)

// Command returns a Command with subcommands for starting and stopping a localkube cluster
func Command(out io.Writer) *cli.Command {
	return &cli.Command{
		Name:        "cluster",
		Description: "Manages localkube Kubernetes development environment",
		Subcommands: []cli.Command{
			{
				Name:        "start",
				Usage:       "spread cluster start [-t <tag>] [ClusterDataDirectory]",
				Description: "Starts the localkube cluster",
				ArgsUsage:   "-t specifies localkube image tag to use, default is latest",
				Action: func(c *cli.Context) {
					// create new Docker client
					ctlr, err := NewControllerFromEnv()
					if err != nil {
						fatal(out, err)
					}

					// create (if needed) and start localkube container
					err = startCluster(ctlr, out, c)
					if err != nil {
						fatal(out, fmt.Errorf("could not start localkube: %v", err))
					}

					// guess which IP the API server will be on
					host, err := identifyHost(ctlr.Endpoint())
					if err != nil {
						fatal(out, err)
					}

					// use default port
					host = fmt.Sprintf("%s:8080", host)

					// setup localkube kubectl context
					currentContext, err := GetCurrentContext()
					if err != nil {
						fatal(out, err)
					}

					// set as current if no CurrentContext is set
					setCurrent := (len(currentContext) == 0)

					err = SetupContext(LocalkubeClusterName, LocalkubeContext, host, setCurrent)
					if err != nil {
						fatal(out, err)
					}

					// display help text messages if context change
					if setCurrent {
						fmt.Fprintf(out, "Created `%s` context and set it as current.\n", LocalkubeContext)
					} else if currentContext != LocalkubeContext {
						fmt.Fprintln(out, SwitchContextInstructions(LocalkubeContext))
					}
				},
			},
			{
				Name:        "stop",
				Usage:       "spread cluster stop [-r]",
				Description: "Stops the localkube cluster",
				ArgsUsage:   "-r removes container",
				Action: func(c *cli.Context) {
					ctlr, err := NewControllerFromEnv()
					if err != nil {
						fatal(out, err)
					}

					stopCluster(ctlr, out, c)
				},
			},
		},
	}
}

// startCluster configures and starts a cluster using command line parameters
func startCluster(ctlr *Controller, out io.Writer, c *cli.Context) error {
	var err error

	// set data directory
	dataDir := c.Args().First()
	if len(dataDir) == 0 {
		dataDir, err = homedir.Expand(DefaultHostDataDir)
		if err != nil {
			return fmt.Errorf("Unable to expand home directory: %v", err)
		}
	}

	// set tag
	tag := c.String("t")
	if len(tag) == 0 {
		tag = LocalkubeDefaultTag
	}

	// check if localkube container exists
	ctrId, running, err := ctlr.OnlyLocalkubeCtr()
	if err != nil {
		if err == ErrNoContainer {
			// if container doesn't exist, create
			ctrId, running, err = ctlr.CreateCtr(LocalkubeContainerName, tag)
			if err != nil {
				return err
			}
		} else {
			// stop for all other errors
			return err
		}
	}

	// start container if not running
	if !running {
		err = ctlr.StartCtr(ctrId, dataDir)
		if err != nil {
			return err
		}
	}
	return nil
}

// stopCluster stops a running cluster
func stopCluster(ctlr *Controller, out io.Writer, c *cli.Context) error {
	remove := c.Bool("r")

	ctrs, err := ctlr.ListLocalkubeCtrs(true)
	if err != nil {
		return err
	}

	for _, ctr := range ctrs {
		ctlr.StopCtr(ctr.ID, remove)
	}
	return nil
}

func identifyHost(endpoint string) (string, error) {
	beginPort := strings.LastIndex(endpoint, ":")
	switch {
	// if using TCP use provided host
	case strings.HasPrefix(endpoint, "tcp://"):
		return endpoint[6:beginPort], nil
	// assuming localhost if Unix
	// TODO: Make this customizable
	case strings.HasPrefix(endpoint, "unix://"):
		return "127.0.0.1", nil
	}
	return "", fmt.Errorf("Could not determine localkube API server from endpoint `%s`", endpoint)
}

func fatal(out io.Writer, err error) {
	fmt.Fprintf(out, "%v\n", err)
	os.Exit(1)
}
