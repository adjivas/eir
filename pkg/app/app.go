package app

import (
	eir_context "github.com/adjivas/eir/internal/context"
	"github.com/adjivas/eir/pkg/factory"
)

type App interface {
	SetLogEnable(enable bool)
	SetLogLevel(level string)
	SetReportCaller(reportCaller bool)

	Start()
	Terminate()

	Context() *eir_context.EIRContext
	Config() *factory.Config
}
