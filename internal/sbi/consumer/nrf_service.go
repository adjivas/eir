package consumer

import (
	"context"
	"strings"
	"sync"
	"time"

	eir_context "github.com/adjivas/eir/internal/context"
	"github.com/adjivas/eir/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/nrf/NFManagement"
	"github.com/pkg/errors"
)

const (
	DELAY_REGISTRATION_NRF = 2 * time.Second
)

type NrfService struct {
	nfMngmntMu sync.RWMutex

	nfMngmntClients map[string]*NFManagement.APIClient
}

func (ns *NrfService) getNFManagementClient(uri string) *NFManagement.APIClient {
	if uri == "" {
		return nil
	}
	ns.nfMngmntMu.RLock()
	client, ok := ns.nfMngmntClients[uri]
	if ok {
		ns.nfMngmntMu.RUnlock()
		return client
	}

	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(uri)
	client = NFManagement.NewAPIClient(configuration)

	ns.nfMngmntMu.RUnlock()
	ns.nfMngmntMu.Lock()
	defer ns.nfMngmntMu.Unlock()
	ns.nfMngmntClients[uri] = client
	return client
}

func (ns *NrfService) buildNFProfile(context *eir_context.EIRContext) (models.NrfNfManagementNfProfile, error) {
	profile := models.NrfNfManagementNfProfile{
		NfInstanceId: context.NfId,
		NfType:       models.NrfNfManagementNfType__5_G_EIR,
		NfStatus:     models.NrfNfManagementNfStatus_REGISTERED,
	}
	if context.RegisterIP.Is6() {
		profile.Ipv6Addresses = []string{context.RegisterIP.String()}
	} else if context.RegisterIP.Is4() {
		profile.Ipv4Addresses = []string{context.RegisterIP.String()}
	}

	var services []models.NrfNfManagementNfService
	for _, nfService := range context.NfService {
		services = append(services, nfService)
	}
	if len(services) > 0 {
		profile.NfServices = services
	}

	return profile, nil
}

func (ns *NrfService) SendRegisterNFInstance(ctx context.Context, nrfUri string) (
	resourceNrfUri string, retrieveNfInstanceId string, err error,
) {
	// Set client and set url
	profile, err := ns.buildNFProfile(eir_context.GetSelf())
	if err != nil {
		logger.ConsumerLog.Errorf("failed to build nrf profile %s", err.Error())
		return "", "", err
	}

	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := ns.getNFManagementClient(nrfUri)

	registerNfInstanceRequest := &NFManagement.RegisterNFInstanceRequest{
		NfInstanceID:             &profile.NfInstanceId,
		NrfNfManagementNfProfile: &profile,
	}
	var res *NFManagement.RegisterNFInstanceResponse
	for {
		select {
		case <-ctx.Done():
			return "", "", errors.Errorf("Context Cancel before RegisterNFInstance")
		default:
		}

		res, err = client.NFInstanceIDDocumentApi.RegisterNFInstance(ctx, registerNfInstanceRequest)
		if err != nil || res == nil {
			logger.ConsumerLog.Errorf("EIR register to NRF Error[%s]", err.Error())
			time.Sleep(DELAY_REGISTRATION_NRF)
			continue
		}

		if res.Location == "" {
			// NFUpdate
			break
		}
		resourceUri := res.Location
		resourceNrfUri, _, _ = strings.Cut(resourceUri, "/nnrf-nfm/")
		retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]

		oauth2 := false

		if res.NrfNfManagementNfProfile.CustomInfo != nil {
			if v, ok := res.NrfNfManagementNfProfile.CustomInfo["oauth2"].(bool); ok == true {
				oauth2 = v
				logger.MainLog.Infoln("OAuth2 setting receive from NRF:", oauth2)
			}
		}

		eir_context.GetSelf().OAuth2Required = oauth2
		if oauth2 && eir_context.GetSelf().NrfCertPem == "" {
			logger.CfgLog.Error("OAuth2 enable but no nrfCertPem provided in config.")
		}
		break
	}
	return resourceNrfUri, retrieveNfInstanceId, err
}

func (ns *NrfService) SendDeregisterNFInstance() (err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	ctx, pd, err := eir_context.GetSelf().GetTokenCtx(models.ServiceName_N5G_EIR_EIC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		logger.ConsumerLog.Errorf("Get token context failed: problem details: %+v", pd)
		return err
	}

	eirSelf := eir_context.GetSelf()

	// Set client and set url
	configuration := NFManagement.NewConfiguration()
	configuration.SetBasePath(eirSelf.NrfUri)
	client := ns.getNFManagementClient(eirSelf.NrfUri)

	if eirSelf.NfId == "" {
		logger.ConsumerLog.Warnf("The EIR haven't a NFId : %+v", eirSelf)
		return nil
	} else {
		deregisterNfInstanceRequest := &NFManagement.DeregisterNFInstanceRequest{
			NfInstanceID: &eirSelf.NfId,
		}
		_, deregisterErr := client.NFInstanceIDDocumentApi.DeregisterNFInstance(ctx, deregisterNfInstanceRequest)
		if deregisterErr != nil {
			return deregisterErr
		}
		return nil
	}
}

func (ns *NrfService) SendSearchNFInstances(nrfUri string,
	param NFDiscovery.SearchNFInstancesRequest,
) (*NFDiscovery.SearchNFInstancesResponse, error) {
	// Set client and set url
	configuration := NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := NFDiscovery.NewAPIClient(configuration)

	ctx, _, err := eir_context.GetSelf().GetTokenCtx(models.ServiceName_NNRF_DISC, models.NrfNfManagementNfType_NRF)
	if err != nil {
		return nil, err
	}

	result, err := client.NFInstancesStoreApi.SearchNFInstances(ctx, &param)

	return result, err
}
