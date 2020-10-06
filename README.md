# QuartzRest

This service provides a web interface for a quartz cluster that is clustered via database tables. 

The main issue it addresses is the need for alerting when triggers go into error state.

In a later version it will also provide full REST accesss to the quartz tables to assist in trouble shooting. 

```
Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
      --port="localhost:8080"  the port to run on
  -P, --postgres-connection=POSTGRES-CONNECTION  
                               Connection string for the postgres database
      --table-prefix="qrtz_"   Prefix of the quartz tables
```

## Endpoints

- `/metrics` is the prometheus metrics endpoint. It has a counter for the errors that reports like this: `quartz_error_counter{sched_name="scheduler",trigger_group="DEFAULT",trigger_name="clientPollerTrigger"} 1667`
This counter increases every 30 seconds when a trigger is in ERROR state.