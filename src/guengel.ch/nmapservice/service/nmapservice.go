package service

import (
	context "context"

	gnms "github.com/RafaelOstertag/grpcnmapservice"
)

// NmapService gRPC Service
type NmapService struct {
}

// Scan given host
func (ns *NmapService) Scan(ctx context.Context, req *gnms.ScanRequest) (*gnms.ScanReply, error) {
	var err error
	var result *Result

	if result, err = Run(req.GetHost(), req.GetPortSpec()); err != nil {
		return nil, err
	}

	return toScanReply(result), nil
}

func toScanReply(nmapResult *Result) *gnms.ScanReply {
	ports := make([]*gnms.ScanReply_Port, len(nmapResult.Ports))

	for i, port := range nmapResult.Ports {
		ports[i] = &gnms.ScanReply_Port{
			Number: int32(port.Number),
			Name: port.Name,
			State: port.State,
		}
	}

	return &gnms.ScanReply{
		State:     nmapResult.State,
		Addresses: nmapResult.Addresses,
		Hostnames: nmapResult.Hostnames,
		Ports: ports,
	}
}
