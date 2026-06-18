package app

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type App struct {
	HTTPServer *http.Server
	GRPCServer *grpc.Server
}

func StartHttpServer(logger *slog.Logger, addr string, mux http.Handler, errCh chan<- error) *http.Server {
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			errCh <- err
		}
	}()

	return srv
}

func StartGRPCServer(
	logger *slog.Logger,
	addr string,
	errCh chan<- error,
	opts []grpc.ServerOption,
	registerService func(grpc.ServiceRegistrar),
) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("failed to listen on gRPC port", "addr", addr, "error", err)
		errCh <- err
		return nil, nil
	}

	srv := grpc.NewServer(opts...)

	if registerService != nil {
		registerService(srv)
	}

	go func() {
		logger.Info("gRPC server started", "addr", addr)
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Error("gRPC server failure", "error", err)
			errCh <- err
		}
	}()

	return srv, lis
}
