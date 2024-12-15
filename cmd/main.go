package main

import (
	"context"
	"fmt"
	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"studentgit.kata.academy/gk/exchanger/cmd/config"
	grpcserver "studentgit.kata.academy/gk/exchanger/internal/grpcServer"
	"studentgit.kata.academy/gk/exchanger/internal/models"
	exchanger "studentgit.kata.academy/gk/exchanger/internal/proto"
	"studentgit.kata.academy/gk/exchanger/internal/repository"
	"studentgit.kata.academy/gk/exchanger/internal/service"
	"syscall"
	"time"
)

var tracer *trace.TracerProvider

func NewTracer(url, name string) error {
	exp, err := jaegerExporter.New(jaegerExporter.WithCollectorEndpoint(jaegerExporter.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tracer = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)
	return nil
}

func InitDb(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.DB_HOST, cfg.DB.DB_USER, cfg.DB.DB_PASS, cfg.DB.DB_NAME, cfg.DB.DB_PORT)

	fmt.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(models.CoinData{})
	if err != nil {
		return nil, fmt.Errorf("failed  to migrate database: %v", err)
	}
	return db, nil
}

func main() {
	logger, err := zap.NewProduction()

	// logger
	log := logger.Sugar()
	if err != nil {
		log.Fatalf("error to init logger : %v", err)
		return
	}

	// config
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Fatalf("error to load config : %v", err)
		return
	}

	err = NewTracer("http://jaeger:14268/api/traces", "server")
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())

	// database
	db, err := InitDb(cfg)
	if err != nil {
		log.Fatalf("error to init db : %v", err)
		return
	}

	// repository
	rep := repository.NewPsqlRepository(db)

	// exchanger
	exch := service.NewGarantexExchange(rep)

	address := "50051"
	lst, err := net.Listen("tcp", ":"+address)
	if err != nil {
		log.Fatalf("error to start listening address %s :%v ", address, err)
	}

	s := grpc.NewServer()

	reflection.Register(s)
	// health server
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	grpcServer := grpcserver.NewExchangerGrpcService(exch, tracer)
	exchanger.RegisterExhangerServer(s, grpcServer)

	// Управление сигналами для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("received signal shutting down gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.GracefulStop()

		select {
		case <-ctx.Done():
			log.Warn("graceful shutdown timed out")
		}

		dbConn, _ := db.DB()
		if err := dbConn.Close(); err != nil {
			log.Errorf("failed to close database connection: %v", err)
		} else {
			log.Info("database connection closed")
		}

		os.Exit(0)
	}()

	if err := s.Serve(lst); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
