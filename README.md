# otelmetricsgin
Open Telemetry Metrics with Golang Gin

There are no worries on this.  The code is terrible, but it works!

It depends on a couple of environment variables and uses gRPC by default.

- `OTEL_EXPORTER_OTLP_ENDPOINT` - example: `http://opentelemetry-collector:4317` or `opentelemetry-collector:4317`.
	If the URL begins with `http://` the library will use an insecure (non-TLS) connection.
	You can also explicitly set `OTEL_EXPORTER_OTLP_INSECURE=true` to force insecure mode.
- `OTEL_SERVICE_NAME`

In main()

```
meterCleanup := otelmetricsgin.InitMeter()
defer meterCleanup(context.Background())
```

I use this after Recovery, Logger and CORS middlewares. If you register it before CORS you may see unknown routes for HTTP OPTIONS requests.

```
if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
	// Apply OpenTelemetry metrics middleware - higher up to catch latency
	router.Use(otelmetricsgin.Middleware())
}
```
