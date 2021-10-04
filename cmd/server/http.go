package main

import (
	"context"
	"net/http"
	"time"
)

func ShutdownOnContext(ctx context.Context, srv *http.Server) {
	<-ctx.Done()
	logger.Info("stopping http server")

	sdctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	srv.Shutdown(sdctx)
}
