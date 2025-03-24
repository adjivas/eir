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
    db_connector := database.NewDbConnector()
	
	return &Processor{
		App:         eir,
		DbConnector: db_connector,
	}
}
