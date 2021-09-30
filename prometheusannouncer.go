package prometheusannouncer

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

type ConfigType uint16

// Configuration Types
const (
	// The announced metrics endpoint is a APP
	APP ConfigType = iota
	// The announced metrics endpoint is a DEVICE
	DEVICE
)

// DialAndAnnounce announces a prometheus metrics endpoint.
// Will run until the context is canceled.
// Connects to the grpc server with insecure. Use Announce and provide your own
// connection to use a secure connection.
// Re-announces the host at interval.
// announceHost and ct are the host definition.
func DialAndAnnounce(
	ctx context.Context,
	grpcServer string,
	interval time.Duration,
	announceHost string,
	ct pb.RegisterRequest_Type,
) error {
	conn, err := grpc.Dial(grpcServer, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not dial grpc server: %v", err)
	}

	return Announce(ctx, conn, interval, announceHost, ct)
}

// Announce announces a prometheus metrics endpoint.
// Will run until the context is canceled.
// Provide an existing connection to the server.
// Re-announces the host at interval.
// announceHost and ct are the host definition.
func Announce(
	ctx context.Context,
	conn *grpc.ClientConn,
	interval time.Duration,
	announceHost string,
	ct pb.RegisterRequest_Type,
) error {
	cli := pb.NewServiceDiscoveryClient(conn)

	req := &pb.RegisterRequest{
		Host:       announceHost,
		ConfigType: ct,
	}

	tick := time.NewTicker(interval)

	for {
		_, err := cli.Announce(ctx, req)
		if err != nil {
			return fmt.Errorf("could not announce service: %v", err)
		}
		select {
		case <-tick.C:
		case <-ctx.Done():
			return nil
		}
	}
}
