// Package localkubectl allows the lifecycle of the localkube container to be controlled
package localkubectl

import (
	"fmt"
	"io"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultHostDataDir is the directory on the host which is mounted for localkube data if no directory is specified
	DefaultHostDataDir = "~/.localkube/data"

	// LocalkubeDefaultTag is the tag to use for the localkube image
	LocalkubeDefaultTag = "latest"
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
					err := startCluster(out, c)
					if err != nil {
						fmt.Fprintf(out, "Could not start localkube: %v", err)
					}

					// TODO: write kubectl context

					// TODO: Set current context to localkube if not set
				},
			},
			{
				Name:        "stop",
				Usage:       "spread cluster stop [-r]",
				Description: "Stops the localkube cluster",
				ArgsUsage:   "-r removes container",
				Action: func(c *cli.Context) {
					stopCluster(out, c)
				},
			},
		},
	}
}

// startCluster configures and starts a cluster using command line parameters
func startCluster(out io.Writer, c *cli.Context) error {
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

	ctlr, err := NewControllerFromEnv()
	if err != nil {
		return err
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
func stopCluster(out io.Writer, c *cli.Context) error {
	remove := c.Bool("r")

	ctlr, err := NewControllerFromEnv()
	if err != nil {
		return err
	}

	ctrs, err := ctlr.ListLocalkubeCtrs(true)
	if err != nil {
		return err
	}

	for _, ctr := range ctrs {
		ctlr.StopCtr(ctr.ID, remove)
	}
	return nil
}
