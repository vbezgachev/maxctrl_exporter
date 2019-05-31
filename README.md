## Overview
This exports the following MaxScale metrics for Prometheus:

- Server connections
- Service session count
- Maxscale instance status
- Event statistics per started thread

## MaxScale requirements
The exporter uses exclusively [MaxScale REST API](https://mariadb.com/kb/en/maxscale-23-rest-api/)

## Installation
1. Install [Golang](https://golang.org/doc/install)
1. Install [dep](https://github.com/golang/dep#installation)
1. Create a new folder in your $GOPATH: `mkdir -p $GOPATH/src/github.com/`
1. Navigate to said folder: `$GOPATH/src/github.com`
1. Clone the repository: `git clone https://github.com/Vetal1977/maxctrl_exporter.git`

## Build
### Manually
1. Change to the project root directory
1. Run `dep ensure --update` to update the dependencies
1. Run `go build` to build the binary for your platform
1. Build linux binary: `GOOS=linux GOARCH=amd64 go build -o bin/linux/maxctrl_exporter`

### With Makefile
1. Run `make build`

## Run locally
1. `cd maxscale_docker`
1. `docker-compose up -d`
1. `docker-compose down` when you are finished

## Test
1. The REST API of MaxScale is accessible at `localhost:8989`. E.g. [http://localhost:8989/v1/servers](http://localhost:8989/v1/servers)
1. The Exporter of MaxScale is accessible at `localhost:8093`. E.g. [http://localhost:8093/metrics](http://localhost:8093/metrics)
