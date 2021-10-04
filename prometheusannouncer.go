package prometheusannouncer

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

// DialAndAnnounce announces a prometheus metrics endpoint.
// Will run until the context is canceled.
// Connects to the grpc server with insecure. Use "Announce" and provide your own
// connection to use a secure connection.
// Re-announces the host at "interval".
// "id" is a unique id (perhaps a uuid) for the client announcement. Supplied by the client.
// ""targets" and "labels" are the target definitions for http_sd.
func DialAndAnnounce(
	ctx context.Context,
	grpcServer string,
	interval time.Duration,
	id string,
	targets []string,
	labels map[string]string,
) error {
	conn, err := grpc.Dial(grpcServer, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not dial grpc server: %v", err)
	}

	return AnnounceWithConnection(ctx, conn, interval, id, targets, labels)
}

// AnnounceWithConnection announces a prometheus metrics endpoint.
// Will run until the context is canceled.
// Provide an existing connection to the server.
// Re-announces the host at "interval".
// "id" is a unique id (perhaps a uuid) for the client announcement. Supplied by the client.
// ""targets" and "labels" are the target definitions for http_sd.
func AnnounceWithConnection(
	ctx context.Context,
	conn *grpc.ClientConn,
	interval time.Duration,
	id string,
	targets []string,
	labels map[string]string,
) error {
	req := Announcement(id, targets, labels)
	return Announce(ctx, conn, interval, req)
}

// Announce announces a prometheus metrics endpoint.
// Will run until the context is canceled.
// Provide an existing connection to the server.
// Re-announces the host at "interval".
// "req" is a RegisterRequest containing an id, targets and labels. Create with "Announcement"
func Announce(ctx context.Context,
	conn *grpc.ClientConn,
	interval time.Duration,
	req *pb.RegisterRequest,
) error {
	cli := pb.NewServiceDiscoveryClient(conn)
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

// Announcement creates a RegisterRequest.
// "id" is a unique id (perhaps a uuid) for the client announcement. Supplied by the client.
// ""targets" and "labels" are the target definitions for http_sd.
func Announcement(id string, targets []string, labels map[string]string) *pb.RegisterRequest {
	ls := []*pb.Label{}
	for l, v := range labels {
		ls = append(ls, &pb.Label{
			Name:  l,
			Value: v,
		})
	}

	return &pb.RegisterRequest{
		UUID:    id,
		Targets: targets,
		Labels:  ls,
	}
}
