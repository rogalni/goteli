package goteli

import (
	"context"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

func setupMetrics(ctx context.Context, con *grpc.ClientConn, resource *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(con))
	if err != nil {
		return nil, err
	}
	mp := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(
			metric.NewPeriodicReader(exporter, metric.WithInterval(15*time.Second)),
		),
	)

	otel.SetMeterProvider(mp)

	err = runtime.Start()
	if err != nil {
		otelzap.S().Error(err)
	}

	return mp, nil
}
