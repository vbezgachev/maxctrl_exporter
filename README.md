# MaxScale Exporter

The Exporter exports the MaxScale metrics for Prometheus:

- Server connections
- Service session count
- MaxScale instance status
- Event statistics per started thread

## MaxScale requirements

The exporter uses exclusively [MaxScale REST API](https://mariadb.com/kb/en/maxscale-23-rest-api/)

## Installation

1. Install [Golang](https://golang.org/doc/install)
1. Create a new folder in your $GOPATH: `mkdir -p $GOPATH/src/github.com/`
1. Navigate to that folder: `$GOPATH/src/github.com`
1. Clone the repository: `git clone https://github.com/Vetal1977/maxctrl_exporter.git`

## Docker

The latest Docker image of `maxctrl_exporter` is available [here](https://github.com/users/vbezgachev/packages/container/package/maxctrl_exporter)

## Build

### Manually

1. Change to the project root directory
1. Run `go build` to build the binary for your platform
1. Build Linux binary: `GOOS=linux GOARCH=amd64 go build -o bin/linux/maxctrl_exporter`

### With Makefile

1. Run `make build`

## Run and test locally

We have prepared a Docker-compose file for a local try. Upon start, you get running MySQL, MaxScale and Exporter containers. Note that the Exported does not use command line arguments rather relies on environment variables:

- MAXSCALE_URL. URL of MaxScale server, default is http://127.0.0.1:8989
- MAXSCALE_USERNAME. MaxScale user name for connection to underlying MySQL database
- MAXSCALE_PASSWORD. MaxScale user password for connection to underlying MySQL database
- MAXSCALE_CA_CERTIFICATE. Certificate to use to verify a secure connection
- MAXSCALE_EXPORTER_PORT. Port that the Exporter expose to provide metrics for Prometheus
- MAXSCALE_MAX_CONNECTIONS. When set will create a static gauge metric with its value, should be set like the [max_connections](https://mariadb.com/kb/en/mariadb-maxscale-25-mariadb-maxscale-configuration-guide/#max_connections) setting in maxscale (this metric can be used to define alerts)

### Run

1. `cd maxscale_docker`
1. `docker-compose up -d`
1. `docker-compose down` when you are finished

### Test

1. The REST API of MaxScale is accessible at `localhost:8989`. E.g. [http://localhost:8989/v1/servers](http://localhost:8989/v1/servers)
1. The Exporter of MaxScale is accessible at `localhost:8093`. E.g. [http://localhost:8093/metrics](http://localhost:8093/metrics)

## Contribution
1. Create a new branch
1. Make changes
1. Run MaxScale and `maxctrl_exporter` locally, e.g. in Docker
1. Check that you can access them and see reasonable output
1. Commit and push the branch
1. Create a PR
1. Request a review from Vitaly Bezgachev, vitaly.bezgachev [the_at_symbol] gmail.com, Kadir Tugan, kadir.tugan [the_at_symbol] gmail.com
