package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"

	pa "github.com/duncanvanzyl/prometheus-announcer"
)

const (
	envPrefix = "httpsd"
)

type Specification struct {
	ID         string        `ignored:"true"`
	Targets    []string      `ignored:"true"`
	Labels     string        `ignored:"true"`
	Interval   time.Duration `default:"1m" desc:"The interval between announcements."`
	GRPCServer string        `default:"localhost:50051" desc:"The GRPC server to announce to.\nIn the form \"host:port\""`
}

func run() error {
	s := &Specification{}

	envconfig.Process(envPrefix, s)

	flag.StringVar(&s.ID, "id", "",
		idUsage,
	)
	flag.StringVar(&s.Labels, "labels", "",
		labelUsage,
	)
	flag.Usage = Usage(envPrefix, s)
	flag.Parse()
	s.Targets = flag.Args()

	if len(s.Targets) < 1 {
		return fmt.Errorf("at least one target is required")
	}

	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	ls, err := processMap(s.Labels)
	if err != nil {
		return fmt.Errorf("could not process labels: %q, %v", s.Labels, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ExitOnSignal(cancel)

	// log.Printf("announcing targets %q with interval %s and labels %q", s.Targets, s.Interval, ls)

	// TODO: check target format
	err = pa.DialAndAnnounce(ctx, s.GRPCServer, s.Interval, s.ID, s.Targets, ls)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
}
