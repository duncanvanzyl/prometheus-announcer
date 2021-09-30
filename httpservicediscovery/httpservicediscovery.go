package httpservicediscovery

import (
	"context"
	"net/http"
	"time"

	"github.com/duncanvanzyl/prometheus-announcer/configstore"
	"github.com/hashicorp/go-hclog"
)

var logger = hclog.Default()

func SetLogger(l hclog.Logger) {
	logger = l.Named("hsd")
}

type HTTPServiceDiscovery interface {
	ShutdownOnContext(context.Context)
	ListenAndServe() error
	// HandleFunc(string, http.HandlerFunc)
	// 	http.Handler
	// 	http.Server
}

type HTTPService struct {
	http.Server
}

func NewHTTPService(host string, f http.Handler) *HTTPService {
	return &HTTPService{
		Server: http.Server{
			Addr:    host,
			Handler: f,
		},
	}
}

func (hs *HTTPService) ShutdownOnContext(ctx context.Context) {

	<-ctx.Done()
	logger.Info("stopping http server")

	sdctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	hs.Server.Shutdown(sdctx)
}

func Handler(cs configstore.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		b, err := cs.JSON()
		if err != nil {
			http.Error(w, "services error", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
