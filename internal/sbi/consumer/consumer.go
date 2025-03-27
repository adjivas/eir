package consumer

import (
	"github.com/adjivas/eir/pkg/app"
	"github.com/free5gc/openapi/nrf/NFManagement"
)

type Consumer struct {
	app.App

	*NrfService
}

func NewConsumer(eir app.App) *Consumer {
	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(eir.Context().NrfUri)
	nrfService := &NrfService{
		nfMngmntClients: make(map[string]*NFManagement.APIClient),
	}

	return &Consumer{
		App:        eir,
		NrfService: nrfService,
	}
}
