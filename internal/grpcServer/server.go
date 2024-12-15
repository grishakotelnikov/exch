package grpcserver

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	exchanger "studentgit.kata.academy/gk/exchanger/internal/proto"
	"studentgit.kata.academy/gk/exchanger/internal/service"
)

type GrpcServer struct {
	exchanger.UnimplementedExhangerServer
	service service.Exchangerer
	tracer  *trace.TracerProvider
}

func NewExchangerGrpcService(exchanger service.Exchangerer, tracer *trace.TracerProvider) *GrpcServer {
	return &GrpcServer{
		service: exchanger,
		tracer:  tracer,
	}
}

func (es *GrpcServer) GetRates(ctx context.Context, req *exchanger.CryptoRequest) (*exchanger.ValueResponce, error) {
	if es.tracer == nil {
		return es.getRatesWithoutTracer(ctx, req)
	}
	ctx, span := es.tracer.Tracer("server").Start(ctx, "GetRates")
	span.SetAttributes(attribute.String("request.currency_pair", req.GetRequest()))
	defer span.End()

	coinData, err := es.service.GetRates(ctx, req.GetRequest())
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(
			attribute.String("error.type", "GetRates"),
			attribute.String("error.message", err.Error()),
		)
		log.Printf("error to get pair = %s currency : %v", req.GetRequest(), err)
		return nil, err
	}
	var res exchanger.ValueResponce
	res.Asks = coinData.AskPrice
	res.Bids = coinData.BidPrice
	res.CurrentTime = coinData.Time
	return &res, nil
}

func (es *GrpcServer) getRatesWithoutTracer(ctx context.Context, req *exchanger.CryptoRequest) (*exchanger.ValueResponce, error) {
	coinData, err := es.service.GetRates(ctx, req.GetRequest())
	if err != nil {
		log.Printf("error to get pair = %s currency : %v", req.GetRequest(), err)
		return nil, err
	}

	return &exchanger.ValueResponce{
		Asks:        coinData.AskPrice,
		Bids:        coinData.BidPrice,
		CurrentTime: coinData.Time,
	}, nil
}
