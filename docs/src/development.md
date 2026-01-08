# Development

## Source code

Source code is written in Go and thus, `go` binaries should be present for compilation and running
tests. You can refer to the [official documentation](https://go.dev/doc/install) for environment
setup.

## Deployment

All supported deployment stacks related files are located into the `/deployment` folder.

### Local

Local deployment is based on [docker-compose](https://docs.docker.com/compose/).

Exposed ports :

| Port | Functionnality |
|------|----------------|
| 8080 | backend        |
| 8081 | swager         |
