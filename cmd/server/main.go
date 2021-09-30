package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/duncanvanzyl/prometheus-announcer/configstore"
	hsd "github.com/duncanvanzyl/prometheus-announcer/httpservicediscovery"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

var logger = hclog.New(&hclog.LoggerOptions{
	Name:              "GPRC-SD",
	Color:             hclog.AutoColor,
	IndependentLevels: false,
})

type app struct {
	pb.UnimplementedServiceDiscoveryServer
	cs         configstore.Config
	httpServer hsd.HTTPServiceDiscovery
}

func run(s *Specification) error {
	logger.SetLevel(hclog.LevelFromString(s.LogLevel))
	configstore.SetLogger(logger)
	hsd.SetLogger(logger)

	logger.Debug("settings", "specification", hclog.Fmt("%+v", *s))

	a := &app{
		cs: configstore.New(s.Lifetime, s.Interval),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hsd.Handler(a.cs))
	mux.HandleFunc("/version", handleVersion())
	a.httpServer = hsd.NewHTTPService(s.HTTPHost, mux)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ExitOnSignal(cancel)

	go a.cs.Run(ctx)

	lis, err := net.Listen("tcp", s.GPRCHost)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	go StopOnContext(ctx, srv)
	pb.RegisterServiceDiscoveryServer(srv, a)

	go func() {
		logger.Info("running grpc server")
		if err := srv.Serve(lis); err != nil {
			logger.Error("failed to serve", "error", err)
			cancel()
		}
	}()

	go a.httpServer.ShutdownOnContext(ctx)

	logger.Info("running http server")
	if err := a.httpServer.ListenAndServe(); err != nil {
		logger.Error("http shutdown error", "error", err)
	}

	return nil
}

type Specification struct {
	PrintVersion bool          `ignored:"true"`
	GPRCHost     string        `default:":50051"`
	LogLevel     string        `default:"info"`
	HTTPHost     string        `default:":8080"`
	Lifetime     time.Duration `default:"2m"`
	Interval     time.Duration `default:"30s"`
}

func main() {
	s := &Specification{}

	flag.BoolVar(&s.PrintVersion, "version", false, "Prints version and exits")
	flag.Parse()

	// print build info
	logger.Info(buildInfo())
	if s.PrintVersion {
		os.Exit(0)
	}

	envconfig.Process("hsd", s)

	err := run(s)
	if err != nil {
		logger.Error("app exited with error", "error", err)
	}
}
