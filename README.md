# hystrix-to-librato

Sends Hystrix/Turbine stream metrics to Librato;

## How it works

You pass down all configs via
[environment variables](/internal/config/config.go) and start up the
binary.

It will create a goroutine for each cluster being monitored, and will
log each report and how many goroutines are alive.

If any HTTP error occurs, it will try to restart itself in 5 seconds.

The reports to Librato are also made in new goroutines, but the same metric
will be sent at most every 5 seconds (to avoid paying too much, since librato
charges by metric sent).

## Metrics sent

- `hystrix.circuit.open`: `1` if circuit open, `0` otherwise.
Source will be `{cluster}.{group}`;
- `hystrix.latency.{lat}`: The executition latency in ms. Source will be
`{cluster}.{group}.{name}`;
