package context

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
)

var eirContext = EIRContext{}

type subsId = string

type EIRServiceType int

const (
	N5G_DR EIRServiceType = iota
)

func Init() {
	eirContext.Name = "eir"
	eirContext.EeSubscriptionIDGenerator = 1
	eirContext.SdmSubscriptionIDGenerator = 1
	eirContext.SubscriptionDataSubscriptionIDGenerator = 1
	eirContext.PolicyDataSubscriptionIDGenerator = 1
	eirContext.SubscriptionDataSubscriptions = make(map[subsId]*models.SubscriptionDataSubscriptions)
	eirContext.PolicyDataSubscriptions = make(map[subsId]*models.PolicyDataSubscription)
	eirContext.InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	serviceName := []models.ServiceName{
		models.ServiceName_N5G_EIR_EIC,
	}
	eirContext.NrfUri = fmt.Sprintf("%s://%s:%d", models.UriScheme_HTTPS, eirContext.RegisterIPv4, 29510)
	initEirContext()

	config := factory.EirConfig
	eirContext.NfService = initNfService(serviceName, config.Info.Version)
}

type EIRContext struct {
	Name                                    string
	UriScheme                               models.UriScheme
	BindingIPv4                             string
	SBIPort                                 int
	NfService                               map[models.ServiceName]models.NrfNfManagementNfService
	RegisterIPv4                            string // IP register to NRF
	HttpIPv6Address                         string
	NfId                                    string
	NrfUri                                  string
	NrfCertPem                              string
	EeSubscriptionIDGenerator               int
	SdmSubscriptionIDGenerator              int
	SubscriptionDataSubscriptionIDGenerator int
	PolicyDataSubscriptionIDGenerator       int
	InfluenceDataSubscriptionIDGenerator    *rand.Rand
	UESubsCollection                        sync.Map // map[ueId]*UESubsData
	UEGroupCollection                       sync.Map // map[ueGroupId]*UEGroupSubsData
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
	InfluenceDataSubscriptions              sync.Map
	appDataInfluDataSubscriptionIdGenerator uint64
	mtx                                     sync.RWMutex
	OAuth2Required                          bool
}

type UESubsData struct {
	EeSubscriptionCollection map[subsId]*EeSubscriptionCollection
	SdmSubscriptions         map[subsId]*models.SdmSubscription
}

type UEGroupSubsData struct {
	EeSubscriptions map[subsId]*models.EeSubscription
}

type EeSubscriptionCollection struct {
	EeSubscriptions      *models.EeSubscription
	AmfSubscriptionInfos []models.AmfSubscriptionInfo
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &EIRContext{}

// Reset EIR Context
func (context *EIRContext) Reset() {
	context.UESubsCollection.Range(func(key, value interface{}) bool {
		context.UESubsCollection.Delete(key)
		return true
	})
	context.UEGroupCollection.Range(func(key, value interface{}) bool {
		context.UEGroupCollection.Delete(key)
		return true
	})
	for key := range context.SubscriptionDataSubscriptions {
		delete(context.SubscriptionDataSubscriptions, key)
	}
	for key := range context.PolicyDataSubscriptions {
		delete(context.PolicyDataSubscriptions, key)
	}
	context.InfluenceDataSubscriptions.Range(func(key, value interface{}) bool {
		context.InfluenceDataSubscriptions.Delete(key)
		return true
	})
	context.EeSubscriptionIDGenerator = 1
	context.SdmSubscriptionIDGenerator = 1
	context.SubscriptionDataSubscriptionIDGenerator = 1
	context.PolicyDataSubscriptionIDGenerator = 1
	context.InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	context.UriScheme = models.UriScheme_HTTPS
	context.Name = "eir"
}

func initEirContext() {
	config := factory.EirConfig
	logger.UtilLog.Infof("eirconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	eirContext.NfId = uuid.New().String()
	eirContext.RegisterIPv4 = factory.EIR_DEFAULT_IPV4 // default localhost
	eirContext.SBIPort = factory.EIR_DEFAULT_PORT_INT  // default port
	if sbi := configuration.Sbi; sbi != nil {
		eirContext.UriScheme = models.UriScheme(sbi.Scheme)
		if sbi.RegisterIPv4 != "" {
			eirContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			eirContext.SBIPort = sbi.Port
		}

		eirContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if eirContext.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			eirContext.BindingIPv4 = sbi.BindingIPv4
			if eirContext.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				eirContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	if configuration.NrfUri != "" {
		eirContext.NrfUri = configuration.NrfUri
	} else {
		logger.UtilLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		eirContext.NrfUri = fmt.Sprintf("%s://%s:%d", eirContext.UriScheme, "127.0.0.1", 29510)
	}
	eirContext.NrfCertPem = configuration.NrfCertPem
}

func initNfService(serviceName []models.ServiceName, version string) (
	nfService map[models.ServiceName]models.NrfNfManagementNfService,
) {
	versionUri := "v" + strings.Split(version, ".")[0]
	nfService = make(map[models.ServiceName]models.NrfNfManagementNfService)
	for idx, name := range serviceName {
		nfService[name] = models.NrfNfManagementNfService{
			ServiceInstanceId: strconv.Itoa(idx),
			ServiceName:       name,
			Versions: []models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          eirContext.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       GetIPv4Uri(),
			IpEndPoints: []models.IpEndPoint{
				{
					Ipv4Address: eirContext.RegisterIPv4,
					Transport:   models.NrfNfManagementTransportProtocol_TCP,
					Port:        int32(eirContext.SBIPort),
				},
			},
		}
	}

	return
}

func GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", eirContext.UriScheme, eirContext.RegisterIPv4, eirContext.SBIPort)
}

func (context *EIRContext) GetIPv4GroupUri(eirServiceType EIRServiceType) string {
	var serviceUri string

	switch eirServiceType {
	case N5G_DR:
		serviceUri = factory.EirDrResUriPrefix
	default:
		serviceUri = ""
	}

	return fmt.Sprintf("%s://%s:%d%s", context.UriScheme, context.RegisterIPv4, context.SBIPort, serviceUri)
}

// Create new EIR context
func GetSelf() *EIRContext {
	return &eirContext
}

func (context *EIRContext) NewAppDataInfluDataSubscriptionID() uint64 {
	context.mtx.Lock()
	defer context.mtx.Unlock()
	context.appDataInfluDataSubscriptionIdGenerator++
	return context.appDataInfluDataSubscriptionIdGenerator
}

func NewInfluenceDataSubscriptionId() string {
	if GetSelf().InfluenceDataSubscriptionIDGenerator == nil {
		GetSelf().InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	}
	return fmt.Sprintf("%08x", GetSelf().InfluenceDataSubscriptionIDGenerator.Uint32())
}

func (c *EIRContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType__5_G_EIR, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}

func (c *EIRContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.UtilLog.Debugf("EIRContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}

	logger.UtilLog.Debugf("EIRContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.NrfCertPem)
}
