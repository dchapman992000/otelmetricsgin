# otelmetricsgin
Open Telemetry Metrics with Golang Gin

There are no worries on this.  The code is terrible, but it works!

It depends on the following environment variables being set and uses grpc instead of http.  In my environment I have a collector listening on both but use grpc

- OTEL_EXPORTER_OTLP_ENDPOINT - "http://opentelemetry-collector:4317" as an example
- OTEL_SERVICE_NAME

In main()

```
meterCleanup := otelmetricsgin.InitMeter()
defer meterCleanup(context.Background())
```

I use this after Recovery, Logger and cors middlewares.  If before CORS, you'll get unknown routes on HTTP OPTIONS

```
if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
	// Apply OpenTelemetry metrics middleware - higher up to catch latency
	router.Use(otelmetricsgin.Middleware())
}
```
