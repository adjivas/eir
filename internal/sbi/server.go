package sbi

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi/models"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/internal/sbi/processor"
	"github.com/adjivas/eir/internal/util"
	"github.com/adjivas/eir/pkg/app"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/free5gc/util/httpwrapper"
	logger_util "github.com/free5gc/util/logger"
)

type Server struct {
	EIR

	httpServer *http.Server
	router     *gin.Engine
}

type EIR interface {
	app.App

	Processor() *processor.Processor
}

func NewServer(eir EIR, tlsKeyLogPath string) *Server {
	s := &Server{
		EIR: eir,
	}

	s.router = newRouter(s)
	server, err := bindRouter(eir, s.router, tlsKeyLogPath)
	s.httpServer = server

	if err != nil {
		logger.SBILog.Errorf("bind Router Error: %+v", err)
		panic("Server initialization failed")
	}

	return s
}

func (s *Server) Run(wg *sync.WaitGroup) {
	logger.SBILog.Info("Starting server...")

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := s.serve()
		if err != http.ErrServerClosed {
			logger.SBILog.Panicf("HTTP server setup failed: %+v", err)
		}
		logger.SBILog.Infof("SBI server (listen on %s) stopped", s.httpServer.Addr)
	}()
}

func (s *Server) Shutdown() {
	s.shutdownHttpServer()
}

func (s *Server) shutdownHttpServer() {
	const shutdownTimeout time.Duration = 2 * time.Second

	if s.httpServer == nil {
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.SBILog.Errorf("HTTP server shutdown failed: %+v", err)
	}
}

func bindRouter(eir app.App, router *gin.Engine, tlsKeyLogPath string) (*http.Server, error) {
	sbiConfig := eir.Config().Configuration.Sbi
	bindAddr := fmt.Sprintf("%s:%d", sbiConfig.BindingIPv4, sbiConfig.Port)

	return httpwrapper.NewHttp2Server(bindAddr, tlsKeyLogPath, router)
}

func newRouter(s *Server) *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)

	dataRepositoryGroup := router.Group(factory.EirDrResUriPrefix)
	dataRepositoryGroup.Use(func(c *gin.Context) {
		util.NewRouterAuthorizationCheck(models.ServiceName_N5G_EIR_EIC).Check(c, s.Context())
	})
	equipementStatusRoutes := s.getEquipementStatusRoutes()
	AddService(dataRepositoryGroup, equipementStatusRoutes)

	return router
}

func (s *Server) unsecureServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) secureServe() error {
	sbiConfig := s.EIR.Config()

	pemPath := sbiConfig.GetCertPemPath()
	if pemPath == "" {
		pemPath = factory.EirDefaultCertPemPath
	}

	keyPath := sbiConfig.GetCertKeyPath()
	if keyPath == "" {
		keyPath = factory.EirDefaultPrivateKeyPath
	}

	return s.httpServer.ListenAndServeTLS(pemPath, keyPath)
}

func (s *Server) serve() error {
	sbiConfig := s.EIR.Config().Configuration.Sbi

	switch sbiConfig.Scheme {
	case "http":
		return s.unsecureServe()
	case "https":
		return s.secureServe()
	default:
		return fmt.Errorf("invalid SBI scheme: %s", sbiConfig.Scheme)
	}
}
