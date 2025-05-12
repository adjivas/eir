package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"

	eir_context "github.com/adjivas/eir/internal/context"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/internal/sbi"
	"github.com/adjivas/eir/internal/sbi/consumer"
	"github.com/adjivas/eir/internal/sbi/processor"
	"github.com/adjivas/eir/pkg/app"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/free5gc/util/mongoapi"
	"github.com/sirupsen/logrus"
)

type EirApp struct {
	cfg    *factory.Config
	eirCtx *eir_context.EIRContext

	ctx    context.Context
	cancel context.CancelFunc

	wg        sync.WaitGroup
	sbiServer *sbi.Server
	processor *processor.Processor
	consumer  *consumer.Consumer
}

var _ app.App = &EirApp{}

func NewApp(ctx context.Context, cfg *factory.Config, tlsKeyLogPath string) (*EirApp, error) {
	eir_context.Init()

	eir := &EirApp{
		cfg:    cfg,
		eirCtx: eir_context.GetSelf(),
		wg:     sync.WaitGroup{},
	}
	eir.ctx, eir.cancel = context.WithCancel(ctx)

	eir.SetLogEnable(cfg.GetLogEnable())
	eir.SetLogLevel(cfg.GetLogLevel())
	eir.SetReportCaller(cfg.GetLogReportCaller())

	processor := processor.NewProcessor(eir)
	eir.processor = processor

	consumer := consumer.NewConsumer(eir)
	eir.consumer = consumer

	eir.sbiServer = sbi.NewServer(eir, tlsKeyLogPath)

	return eir, nil
}

func (a *EirApp) Consumer() *consumer.Consumer {
	return a.consumer
}

func (a *EirApp) Processor() *processor.Processor {
	return a.processor
}

func (a *EirApp) Config() *factory.Config {
	return a.cfg
}

func (a *EirApp) Context() *eir_context.EIRContext {
	return a.eirCtx
}

func (a *EirApp) SetLogEnable(enable bool) {
	logger.MainLog.Infof("Log enable is set to [%v]", enable)
	if enable && logger.Log.Out == os.Stderr {
		return
	} else if !enable && logger.Log.Out == io.Discard {
		return
	}
	a.cfg.SetLogEnable(enable)
	if enable {
		logger.Log.SetOutput(os.Stderr)
	} else {
		logger.Log.SetOutput(io.Discard)
	}
}

func (a *EirApp) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.MainLog.Warnf("Log level [%s] is invalid", level)
		return
	}
	logger.MainLog.Infof("Log level is set to [%s]", level)
	if lvl == logger.Log.GetLevel() {
		return
	}
	a.cfg.SetLogLevel(level)
	logger.Log.SetLevel(lvl)
}

func (a *EirApp) SetReportCaller(reportCaller bool) {
	logger.MainLog.Infof("Report Caller is set to [%v]", reportCaller)
	if reportCaller == logger.Log.ReportCaller {
		return
	}
	a.cfg.SetLogReportCaller(reportCaller)
	logger.Log.SetReportCaller(reportCaller)
}

func (u *EirApp) registerToNrf(ctx context.Context) error {
	eirContext := u.eirCtx

	nrfUri, nfId, err := u.consumer.SendRegisterNFInstance(ctx, eirContext.NrfUri)
	if err != nil {
		return fmt.Errorf("send register NFInstance error[%s]", err.Error())
	}
	eirContext.NrfUri = nrfUri
	eirContext.NfId = nfId

	return nil
}

func (a *EirApp) deregisterFromNrf() {
	err := a.consumer.SendDeregisterNFInstance()
	if err != nil {
		switch apiErr := err.(type) {
		case openapi.GenericOpenAPIError:
			switch errModel := apiErr.Model().(type) {
			case NFManagement.DeregisterNFInstanceError:
				pd := &errModel.ProblemDetails
				logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", pd)
			case error:
				logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
			}
		case error:
			logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
		}
		logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
	}

	logger.InitLog.Infof("Deregister from NRF successfully")
}

func (a *EirApp) Start() {
	// get config file info
	config := factory.EirConfig
	mongodb := config.Configuration.Mongodb

	// Connect to MongoDB
	if err := mongoapi.SetMongoDB(mongodb.Name, mongodb.Url); err != nil {
		logger.InitLog.Errorf("EIR start set MongoDB error: %+v", err)
		return
	}

	// Register to Nrf
	err := a.registerToNrf(a.ctx)
	if err != nil {
		logger.InitLog.Errorf("register to NRF failed: %v", err)
	} else {
		logger.InitLog.Infof("register to NRF successfully")
	}

	// Graceful deregister when panic
	defer func() {
		if p := recover(); p != nil {
			logger.InitLog.Errorf("panic: %v\n%s", p, string(debug.Stack()))
			a.deregisterFromNrf()
		}
	}()

	logger.InitLog.Infoln("Server started")

	a.wg.Add(1)
	go a.listenShutdown(a.ctx)

	a.sbiServer.Run(&a.wg)
	a.WaitRoutineStopped()
}

func (a *EirApp) listenShutdown(ctx context.Context) {
	defer a.wg.Done()

	<-ctx.Done()
	a.terminateProcedure()
}

func (a *EirApp) Terminate() {
	a.cancel()
}

func (a *EirApp) terminateProcedure() {
	logger.MainLog.Infof("Terminating EIR...")
	a.CallServerStop()
	a.deregisterFromNrf()
}

func (a *EirApp) CallServerStop() {
	if a.sbiServer != nil {
		a.sbiServer.Shutdown()
	}
}

func (a *EirApp) WaitRoutineStopped() {
	a.wg.Wait()
	logger.MainLog.Infof("EIR terminated")
}
