package main

import (
	"context"

	"google.golang.org/grpc"

	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

func StopOnContext(ctx context.Context, srv *grpc.Server) {
	<-ctx.Done()
	logger.Info("stopping grpc server")
	srv.Stop()
}

func (a *app) Announce(ctx context.Context, rr *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id := rr.GetUUID()
	ts := rr.GetTargets()
	ls := make(map[string]string)
	for _, l := range rr.GetLabels() {
		ls[l.Name] = l.Value
	}

	a.cs.AddTarget(id, ts, ls)

	return &pb.RegisterResponse{}, nil
}
