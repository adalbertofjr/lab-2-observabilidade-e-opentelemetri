package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/cmd/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/api"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/gateway"
	"github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/infra/web"

	usecase "github.com/adalbertofjr/lab-1-go-weather-cloud-run/internal/usecase/weather"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	shutdown, err := initProvider(configs.OtelEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown TracerProvider: %v", err)
		}
	}()

	tracer := otel.Tracer("serviceB-tracer")

	startServer(configs, tracer)
}

func startServer(configs *configs.Conf, tracer trace.Tracer) {
	weatherGateway := gateway.NewWeatherAPI(configs.WeatherAPIKey, tracer)
	weatherUseCase := usecase.NewWeatherUseCase(weatherGateway, tracer)
	weatherHandler := api.NewWeatherHandler(weatherUseCase, tracer)
	healthHandler := api.NewHealthCheck()

	webserver := web.NewWebServer(configs.WebServerPort)
	webserver.AddHandler("/", weatherHandler.GetWeather)
	webserver.AddHandler("/health", healthHandler.HealthCheck)
	webserver.Start()
}

func initProvider(otelEndpoint string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ServiceB"),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if otelEndpoint == "" {
		otelEndpoint = "localhost:4317"
	}

	conn, err := grpc.NewClient(
		otelEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC connection to collector: %w", err)
	}

	tracerExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("Failed to create the collector trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(tracerExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}
