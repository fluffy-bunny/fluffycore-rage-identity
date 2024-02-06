package mystream

import (
	"context"
	"io"
	"math/rand"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_config "github.com/fluffy-bunny/fluffycore-starterkit/internal/contracts/config"
	fluffycore_contracts_somedisposable "github.com/fluffy-bunny/fluffycore-starterkit/internal/contracts/somedisposable"
	proto_helloworld "github.com/fluffy-bunny/fluffycore-starterkit/proto/helloworld"
	endpoint "github.com/fluffy-bunny/fluffycore/contracts/endpoint"
	grpc_gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/rs/zerolog/log"
	grpc "google.golang.org/grpc"
)

type (
	service struct {
		proto_helloworld.MyStreamServiceFluffyCoreServer

		config               *contracts_config.Config
		scopedSomeDisposable fluffycore_contracts_somedisposable.IScopedSomeDisposable
	}
)

var stemService = (*service)(nil)

func (s *service) Ctor(config *contracts_config.Config,
	scopedSomeDisposable fluffycore_contracts_somedisposable.IScopedSomeDisposable) proto_helloworld.IFluffyCoreMyStreamServiceServer {
	return &service{
		config:               config,
		scopedSomeDisposable: scopedSomeDisposable,
	}
}

func init() {
	var _ proto_helloworld.IFluffyCoreMyStreamServiceServer = (*service)(nil)
	var _ endpoint.IEndpointRegistration = (*service)(nil)

}

func (s *service) RegisterHandler(gwmux *grpc_gateway_runtime.ServeMux, conn *grpc.ClientConn) {
	proto_helloworld.RegisterMyStreamServiceHandler(context.Background(), gwmux, conn)
}

func AddMyStreamService(builder di.ContainerBuilder) {
	proto_helloworld.AddMyStreamServiceServerWithExternalRegistration(builder,
		stemService.Ctor, func() endpoint.IEndpointRegistration {
			return &service{}
		})
}

func (s *service) RequestPoints(request *proto_helloworld.PointsRequest, stream proto_helloworld.MyStreamService_RequestPointsServer) error {

	numPoints := len(request.Points)
	for i := 0; i < numPoints; i++ {
		for j := 0; j < 10; j++ {
			err := stream.Send(&proto_helloworld.Point{
				Latitude:  rand.Int31n(100),
				Longitude: rand.Int31n(100),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (s *service) StreamPoints(stream proto_helloworld.MyStreamService_StreamPointsServer) error {
	var pointCount int32
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		log.Info().Interface("point", point).Msg("StreamPoints")
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&proto_helloworld.RouteSummary{
				PointCount: pointCount,

				ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++

	}
}
