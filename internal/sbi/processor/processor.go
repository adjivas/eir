package processor

import (
	"github.com/adjivas/eir/internal/database"
	"github.com/adjivas/eir/pkg/app"
	"github.com/adjivas/eir/internal/logger"
)

type Processor struct {
	app.App
	database.DbConnector
}

func NewProcessor(eir app.App) *Processor {
	db_connector, err := database.NewDbConnector()
	if err != nil {
		logger.InitLog.Errorf("DB Connector error: %+v", err)
		panic("processor initialization failed")
	}
	has, err := db_connector.HasEquipementStatus()
	if err != nil {
		logger.InitLog.Errorf("Has EquipementStatus error: %+v", err)
		panic("processor initialization failed")
	}
	if has == false {
		logger.InitLog.Infof("Haven't found EquipementStatus, will create one")
		if err = db_connector.CreateEquipementStatus(); err != nil {
			logger.InitLog.Errorf("Create EquipementStatus error: %+v", err)
			panic("processor initialization failed")
		}
	}
	
	return &Processor{
		App:         eir,
		DbConnector: db_connector,
	}
}
