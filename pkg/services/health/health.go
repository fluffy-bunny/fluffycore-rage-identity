package health

import (
	"context"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_health "github.com/fluffy-bunny/fluffycore/contracts/health"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
)

type service struct{
	grpc_health.UnimplementedHealthServer
}

var _ fluffycore_contracts_health.IHealthServer = (*service)(nil)
func (s *service) Check(context.Context, *grpc_health.HealthCheckRequest) (*grpc_health.HealthCheckResponse, error) {
	return &grpc_health.HealthCheckResponse{
		Status: grpc_health.HealthCheckResponse_SERVING,
	}, nil
}
func (s *service) Watch(*grpc_health.HealthCheckRequest, grpc_health.Health_WatchServer) error {
	return nil
}
func AddHealthService(cb di.ContainerBuilder) {
	di.AddSingleton[fluffycore_contracts_health.IHealthServer](cb, func() fluffycore_contracts_health.IHealthServer {
		return &service{}
	})
}
