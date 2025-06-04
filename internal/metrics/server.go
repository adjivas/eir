package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
	cfg        *factory.Config
}

// Initializes a new HTTP server instance and associate the prometheus handler to it
func NewServer(cfg *factory.Config, tlsKeyLogPath string) (*Server, error) {
	mux := http.NewServeMux()
	reg := Init(cfg)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	bindAddr := cfg.GetMetricsBindingAddr()
	logger.MetricsLog.Infof("Binding addr: [%s]", bindAddr)

	httpServer, err := httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, mux)

	if err != nil {
		logger.MetricsLog.Errorf("Initialize HTTP server failed: %v", err)
		return nil, err
	}

	s := &Server{
		httpServer: httpServer,
		cfg:        cfg,
	}

	return s, nil
}

// Configure the server to handle http requests
func (s *Server) ListenAndServe() {
	logger.MetricsLog.Infof("Starting HTTP server on %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()

	if err != nil {
		logger.MetricsLog.Errorf("Metric server error: %v", err)
	}
}

// Configure the server to handle https requests
func (s *Server) ListenAndServeTLS() {
	tlsKeyPath, tlsPemPath := s.cfg.GetCertKeyPath(), s.cfg.GetCertPemPath()

	err := s.httpServer.ListenAndServeTLS(tlsKeyPath, tlsPemPath)

	if err != nil {
		logger.MetricsLog.Errorf("Metric server error: %v", err)
	}
}

func (s *Server) startServer(wg *sync.WaitGroup) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MetricsLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
		wg.Done()
	}()

	var err error

	logger.MetricsLog.Infof("Start Metrics server (listen on %s)", s.httpServer.Addr)
	scheme := s.cfg.GetMetricsScheme()

	if scheme == "http" {
		err = s.httpServer.ListenAndServe()
	} else if scheme == "https" {
		tlsKeyPath, tlsPemPath := s.cfg.GetCertKeyPath(), s.cfg.GetCertPemPath()
		err = s.httpServer.ListenAndServeTLS(tlsKeyPath, tlsPemPath)
	} else {
		err = fmt.Errorf("no support this scheme[%s]", scheme)
	}

	if err != nil && err != http.ErrServerClosed {
		logger.MetricsLog.Errorf("Metrics server error: %v", err)
	}
	logger.MetricsLog.Warnf("Metrics server (listen on %s) stopped", s.httpServer.Addr)

}

func (s *Server) Run(cfg *factory.Config, wg *sync.WaitGroup) {
	wg.Add(1)
	go s.startServer(wg)
}

func (s *Server) Stop() {
	const defaultShutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer != nil {
		logger.MetricsLog.Infof("Stop Metrics server (listen on %s)", s.httpServer.Addr)
		toCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(toCtx); err != nil {
			logger.MetricsLog.Errorf("Could not close Metrics server: %#v", err)
		}
	}
}
