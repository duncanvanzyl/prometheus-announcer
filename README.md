# Prometheus Announcer

A gprc server that allows applications to announce themselves to prometheus 
using [prometheus http service discovery](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config).

Connect using grpc and Announce the service to have it presented using the prometheus http service
discovery. The app to be discovered is responsible for reannouncing itself before the announcement 
lifetime has expired.

Example Announce Function:
```go
import "github.com/duncanvanzyl/prometheus-announcer/pb"
⋮
go func(
		ctx context.Context,
		grpcServer string,
		interval time.Duration,
		announceHost string,
		ct pb.RegisterRequest_Type,
	) {
		conn, err := grpc.Dial(grpcServer, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not dial grpc server: %v", err)
		}

		cli := pb.NewServiceDiscoveryClient(conn)

		req := &pb.RegisterRequest{
			Host:       announceHost,
			ConfigType: ct,
		}

		tick := time.NewTicker(interval)

		for {
			_, err = cli.Announce(ctx, req)
			if err != nil {
				log.Fatalf("could not announce service: %v", err)
			}
			select {
			case <-tick.C:
			case <-ctx.Done():

				return
			}
		}
	}(ctx, grpcServer, 5*time.Minute, announceHost, pb.RegisterRequest_PIAPP)
```

Or just use the provided Announce functions:
```go
import (
	pa "github.com/duncanvanzyl/prometheus-announcer"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)
⋮
// Note: It is probably worth doing something with the returned error value...
go pa.DialAndAnnounce(ctx, grpcServer, 1*time.Minute, announceHost, pb.RegisterRequest_APP)
```