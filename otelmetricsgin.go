package otelmetricsgin

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

var requestCount metric.Int64Counter
var requestDuration metric.Float64Histogram

func InitMeter() func(context.Context) error {

	//TODO: Only do cleanup if we're using OTLP
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return func(ctx context.Context) error {
			log.Print("nil cleanup function - success if this is without OTEL!")
			return nil
		}
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
	)

	if err != nil {
		log.Fatalf("Failed to create exporter: %v", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatalf("Could not set resources: %v", err)
	}

	// Register the exporter with an SDK via a periodic reader.
	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(res),
		metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter)),
	)

	otel.SetMeterProvider(provider)

	meter := otel.Meter("gin-opentelemetry")

	requestCount, err = meter.Int64Counter(
		"http_server_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)

	if err != nil {
		log.Fatalf("failed to create requestCount instrument: %v", err)
	}

	requestDuration, err = meter.Float64Histogram(
		"http_server_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		// Default buckets are much more than this
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.05, 0.5, 1, 5),
	)
	if err != nil {
		log.Fatalf("failed to create requestDuration instrument: %v", err)
	}

	return exporter.Shutdown
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Let the other handlers process now that the time has started
		c.Next()

		duration := time.Since(start).Seconds()

		route := c.FullPath()
		if len(route) <= 0 {
			route = "nonconfigured_route"
		}

		// Group the codes
		// TODO: make this a config option later?
		code := int(c.Writer.Status()/100) * 100

		requestCount.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.Int("http.status_code", code),
				attribute.String("http.path", route),
				attribute.String("instance", os.Getenv("HOSTNAME")),
			),
		)

		requestDuration.Record(c.Request.Context(), duration,
			metric.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.Int("http.status_code", code),
				attribute.String("http.path", route),
				attribute.String("instance", os.Getenv("HOSTNAME")),
			),
		)
	}
}
