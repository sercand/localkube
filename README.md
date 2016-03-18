# localkube
[![Build Status](https://travis-ci.org/redspread/localkube.svg?branch=master)](https://travis-ci.org/redspread/localkube)

`localkube` provides a Kubernetes cluster configured to run locally and optimized for rapid development. 

The environment is intended to be used with Kubernetes development tool [spread](https://github.com/redspread/spread). It is similar to [monokube](https://github.com/polvi/monokube).

### Developing

This will compile localkube executable, build an image with it inside, and run a container with the build image.

The `docker` command should be setup for the Docker daemon you want the run the cluster with.

**Linux**
```bash
make run-image
```

**OS X/Windows**
```bash
make docker-build run-image
```

The apiserver will run on port 8080 of the Docker host.