# QuartzRest
![CI](https://github.com/UnikumAB/quartzRest/workflows/CI/badge.svg)

This service provides a web interface for a quartz cluster that is clustered via database tables. 

The main issue it addresses is the need for alerting when triggers go into error state.

In a later version it will also provide full REST accesss to the quartz tables to assist in trouble shooting. 
##Usage
```
QuartzRest is a simple tool for accessing quartz data tables.

Usage:
  ./quartzRestServer [flags]

Flags:
  -b, --bind string                  the host:port to bind to (default "localhost:8080")
      --config string                config file (default is $HOME/.cobra.yaml)
  -h, --help                         help for ./quartzRestServer
  -P, --postgres-connection string   Connection string for the postgres database
  -p, --table-prefix string          Prefix of the quartz tables (default "qrtz_")
      --viper                        use Viper for configuration (default true)
```
## Environment variables

| variable Name| Description |
|:----|:----|
| QUARTZ_SERVER_BIND | same as --bind | 
| QUARTZ_SERVER_POSTGRES_CONNECTION| same as --postgres-connection |
| QUARTZ_SERVER_TABLE_PREFIX | same as --table-prefix |

## Endpoints

- `/metrics` is the prometheus metrics endpoint. It has a counter for the errors that reports like this: `quartz_error_counter{sched_name="scheduler",trigger_group="DEFAULT",trigger_name="clientPollerTrigger"} 1667`
This counter increases every 30 seconds when a trigger is in ERROR state.

## Docker image
run 
```
$ docker run unikum/quartz-rest-server
```
