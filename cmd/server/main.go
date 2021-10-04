package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/duncanvanzyl/prometheus-announcer/configstore"
	"github.com/duncanvanzyl/prometheus-announcer/pb"
)

var logger = hclog.New(&hclog.LoggerOptions{
	Name:              "GPRC-SD",
	Color:             hclog.AutoColor,
	IndependentLevels: false,
})

type app struct {
	pb.UnimplementedServiceDiscoveryServer
	cs configstore.ConfigurationStore
	// httpServer hsd.HTTPServiceDiscovery
	router *mux.Router
}

func run(s *Specification) error {
	logger.SetLevel(hclog.LevelFromString(s.LogLevel))
	configstore.SetLogger(logger)

	logger.Debug("settings", "specification", hclog.Fmt("%+v", *s))

	a := &app{
		cs:     configstore.New(s.Lifetime, s.Interval),
		router: mux.NewRouter(),
	}

	a.routes(s.WithREST)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ExitOnSignal(cancel)

	go a.cs.Run(ctx)

	// Create and run the GRPC Server
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

	// Create and run the HTTP server
	hSrv := &http.Server{Addr: s.HTTPHost, Handler: a.router}
	go ShutdownOnContext(ctx, hSrv)
	// go a.httpServer.ShutdownOnContext(ctx)
	logger.Info("running http server")
	if err := hSrv.ListenAndServe(); err != nil {
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
	WithREST     bool          `default:"true"`
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
