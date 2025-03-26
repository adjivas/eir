package processor

import (
	"github.com/adjivas/eir/internal/database"
	"github.com/adjivas/eir/pkg/app"
)

type Processor struct {
	app.App
	database.DbConnector
}

func NewProcessor(eir app.App) *Processor {
	return &Processor{
		App:         eir,
		DbConnector: database.NewDbConnector(eir.Config().Configuration.DbConnectorType),
	}
}
