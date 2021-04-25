package service

import (
	context "context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gnms "github.com/RafaelOstertag/grpcnmapservice"
)

// NmapService gRPC Service
type NmapService struct {
}

// Scan given host
func (ns *NmapService) Scan(ctx context.Context, req *gnms.ScanRequest) (*gnms.ScanReply, error) {
	var err error
	var result *Result

	log.Printf("Scan host '%s' with portspec '%s'", req.GetHost(), req.GetPortSpec())

	if result, err = Run(req.GetHost(), req.GetPortSpec()); err != nil {
		if _, ok := err.(HostSpecError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if _, ok := err.(PortSpecError); ok {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if _, ok := err.(ScannerError); ok {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return nil, status.Error(codes.Unknown, err.Error())
	}

	return toScanReply(result), nil
}

func toScanReply(nmapResult *Result) *gnms.ScanReply {
	ports := make([]*gnms.ScanReply_Port, len(nmapResult.Ports))

	for i, port := range nmapResult.Ports {
		ports[i] = &gnms.ScanReply_Port{
			Number: int32(port.Number),
			Name:   port.Name,
			State:  port.State,
		}
	}

	return &gnms.ScanReply{
		State:     nmapResult.State,
		Addresses: nmapResult.Addresses,
		Hostnames: nmapResult.Hostnames,
		Ports:     ports,
	}
}
