package goteli

import (
	"context"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Opts struct {
	ServiceName           string
	LogLevel              string
	IsJsonLogging         bool
	IsTracingEnabled      bool
	IsMetricsEnabled      bool
	GrpcCollectorEndpoint string
}

func NewDefaultOpts() Opts {
	return Opts{
		ServiceName:           "goteli",
		LogLevel:              "INFO",
		IsJsonLogging:         true,
		IsTracingEnabled:      true,
		IsMetricsEnabled:      true,
		GrpcCollectorEndpoint: "localhost:4317",
	}
}

type Goteli struct {
	cleanupTracerFunc         func(context.Context) error
	cleanupMeterFunc          func(context.Context) error
	cleanupGrpcConnectionFunc func() error
	syncLoggerFunc            func() error
}

func New(ctx context.Context, opts Opts) func(context.Context) {

	goteli := &Goteli{}

	sync := setupLogger(opts.LogLevel, opts.IsJsonLogging, opts.ServiceName)
	goteli.syncLoggerFunc = sync
	if !opts.IsTracingEnabled && !opts.IsMetricsEnabled {
		return goteli.cleanup
	}

	r, err := otelResource(ctx, opts.ServiceName)
	if err != nil {
		otelzap.S().Error(err)
	}
	conn := otelGrpcCon(ctx, opts.GrpcCollectorEndpoint)
	goteli.cleanupGrpcConnectionFunc = conn.Close

	if opts.IsTracingEnabled {
		tp, err := setupTracing(ctx, conn, r)
		if err != nil {
			otelzap.S().Error(err)
		}
		goteli.cleanupTracerFunc = tp.Shutdown
	}

	if opts.IsMetricsEnabled {
		mp, err := setupMetrics(ctx, conn, r)
		if err != nil {
			otelzap.S().Error(err)
		}
		goteli.cleanupMeterFunc = mp.Shutdown
	}

	return goteli.cleanup
}

func otelGrpcCon(ctx context.Context, endpoint string) *grpc.ClientConn {
	con, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		otelzap.S().Error(err)
	}
	return con
}

// Returns a new OpenTelemetry resource describing this application.
func otelResource(ctx context.Context, serviceName string) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
}

func (o *Goteli) cleanup(ctx context.Context) {
	if o.cleanupTracerFunc != nil {
		otelzap.Ctx(ctx).Debug("Cleanup Tracer")
		err := o.cleanupTracerFunc(ctx)
		if err != nil {
			otelzap.Ctx(ctx).Warn("Cleanup meter failed")
		}
	}

	if o.cleanupMeterFunc != nil {
		otelzap.Ctx(ctx).Debug("Cleanup Meter")
		err := o.cleanupMeterFunc(ctx)
		if err != nil {
			otelzap.Ctx(ctx).Warn("Cleanup meter failed")
		}
	}

	if o.cleanupGrpcConnectionFunc != nil {
		otelzap.Ctx(ctx).Debug("Cleanup grpc connection")
		err := o.cleanupGrpcConnectionFunc()
		if err != nil {
			otelzap.Ctx(ctx).Warn("Cleanup meter failed")
		}
	}

	if o.syncLoggerFunc != nil {
		otelzap.Ctx(ctx).Debug("Cleanup otel logger")
		err := o.syncLoggerFunc()
		if err != nil {
			otelzap.Ctx(ctx).Warn("Cleanup meter failed")
		}
	}
}
