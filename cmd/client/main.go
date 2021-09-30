package main

import (
	"context"
	"time"

	pa "github.com/duncanvanzyl/prometheus-announcer"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

const (
	grpcServer   = "localhost:50051"
	announceHost = "localhost:2112"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Note: It is probably worth doing something with the returned error value...
	go pa.DialAndAnnounce(ctx, grpcServer, 1*time.Minute, announceHost, pb.RegisterRequest_APP)

	<-ctx.Done()
}
