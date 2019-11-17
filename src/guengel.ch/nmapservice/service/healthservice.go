package service

import (
	context "context"

	gnms "github.com/RafaelOstertag/grpcnmapservice"
)

const (
	grpcServiceName = "ch.guengel.nmapservice.Nmap"
)

var (
	serviceOk = gnms.HealthCheckResponse{
		Status: gnms.HealthCheckResponse_SERVING,
	}
	serviceUnknown = gnms.HealthCheckResponse{
		Status: gnms.HealthCheckResponse_UNKNOWN,
	}
)

// HealthService implementing GRPC Health Checking Protocol
type HealthService struct {
	// Health Channel to receive health status updates
	Health chan gnms.HealthCheckResponse_ServingStatus
}

func (hs *HealthService) Check(ctx context.Context, in *gnms.HealthCheckRequest) (*gnms.HealthCheckResponse, error) {
	requestedServiceName := in.GetService()
	if requestedServiceName == "" || requestedServiceName == grpcServiceName {
		return &serviceOk, nil
	}

	return &serviceUnknown, nil
}

func (hs *HealthService) Watch(in *gnms.HealthCheckRequest, stream gnms.Health_WatchServer) error {
	requestedServiceName := in.GetService()

	if requestedServiceName != "" && requestedServiceName != serviceName {
		if err := stream.Send(&serviceUnknown); err != nil {
			return err
		}
		return nil
	}

	// Send the first status
	if err := stream.Send(&serviceOk); err != nil {
		return err
	}

	// Indefenitely listen to status changes
	for {
		status := <-hs.Health
		if err := stream.Send(&gnms.HealthCheckResponse{Status: status}); err != nil {
			return err
		}
	}
}
