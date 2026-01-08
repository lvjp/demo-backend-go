# Configuration

Default path as defined in the source code :

```go{{#include ../../internal/app/config/config.go:default_config_path}}```

## Example

This is the actual configuration file used for the local deployment.

```yaml
{{#include ../../deployments/local/config.yaml}}
```

## Syntax

`server.listen_address`:
  : Address used to listen for incoming HTTP requests with [fiber/App.Listen][fiber/App.Listen].

`log.level`:
  : Level used for logging.  
    Valid values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`

`log.format`:
  : Logging output format.  
    Valid values :
    - `json`: JSON formatted output
    - `console`: Shiny debugging colored output for console

`database.driver_name`:
  : Passed as first argument of [database/sql.Open][database/sql.Open] on database connection.

`database.data_source_name`:
  : Passed as first argument of [database/sql.Open][database/sql.Open] on database connection.

[fiber/App.Listen]: https://pkg.go.dev/github.com/gofiber/fiber/v2#App.Listen
[database/sql.Open]: https://pkg.go.dev/database/sql#Open
