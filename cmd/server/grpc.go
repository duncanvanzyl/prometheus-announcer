package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pa "github.com/duncanvanzyl/prometheus-announcer"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

func StopOnContext(ctx context.Context, srv *grpc.Server) {
	<-ctx.Done()
	// logger.Info("stopping grpc server")
	logger.Info("stopping grpc server")
	// srv.GracefulStop()
	srv.Stop()
}

func (a *app) Announce(ctx context.Context, rr *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	host := rr.GetHost()
	var ct pa.ConfigType
	switch rr.GetConfigType() {
	case pb.RegisterRequest_APP:
		ct = pa.APP
	case pb.RegisterRequest_DEVICE:
		ct = pa.DEVICE
	default:
		return nil, status.Errorf(codes.Unimplemented, "invalid config type: %v", rr.GetConfigType())
	}

	a.cs.AddTarget(host, ct)

	return &pb.RegisterResponse{}, nil
}
