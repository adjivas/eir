package context

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/adjivas/eir/internal/logger"
	"github.com/adjivas/eir/pkg/factory"
	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
)

var eirContext = EIRContext{}

func Init() {
	eirContext.Name = "eir"

	serviceName := []models.ServiceName{
		models.ServiceName_N5G_EIR_EIC,
	}

	eirContext.NrfUri = GetIPUri()
	initEirContext()

	config := factory.EirConfig
	eirContext.NfService = initNfService(serviceName, config.Info.Version)
}

type EIRContext struct {
	Name            string
	UriScheme       models.UriScheme
	RegisterIP      netip.Addr // IP register to NRF
	BindingIP       netip.Addr
	SBIPort         int
	DefaultStatus   string
	NfService       map[models.ServiceName]models.NrfNfManagementNfService
	NfId            string
	NrfUri          string
	NrfCertPem      string
	OAuth2Required  bool
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &EIRContext{}

func initEirContext() {
	config := factory.EirConfig
	logger.UtilLog.Infof("eirconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	configuration := config.Configuration
	eirContext.NfId = uuid.New().String()
	sbi := configuration.Sbi

	eirContext.SBIPort = sbi.Port                       // default port
	eirContext.UriScheme = models.UriScheme(sbi.Scheme) // default localhost

	if bindingIP := os.Getenv(sbi.BindingIP); bindingIP != "" {
		logger.UtilLog.Info("Parsing BindingIP address from ENV Variable.")
		sbi.BindingIP = bindingIP
	}
	if registerIP := os.Getenv(sbi.RegisterIP); registerIP != "" {
		logger.UtilLog.Info("Parsing RegisterIP address from ENV Variable.")
		sbi.RegisterIP = registerIP
	}

	eirContext.BindingIP = resolveIP(sbi.BindingIP)
	eirContext.RegisterIP = resolveIP(sbi.RegisterIP)

	eirContext.NrfUri = configuration.NrfUri
	eirContext.NrfCertPem = configuration.NrfCertPem

	if defaultStatus := configuration.DefaultStatus; defaultStatus != "" {
		eirContext.DefaultStatus = defaultStatus
	}

	fmt.Println("eir context = ", &eirContext)
}

func resolveIP(ip string) netip.Addr {
	resolvedIPs, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", ip)
	if err != nil {
		logger.InitLog.Errorf("Lookup failed with %s: %+v", ip, err)
	}
	resolvedIP := resolvedIPs[0].Unmap()
	if resolvedIP := resolvedIP.String(); resolvedIP != ip {
		logger.UtilLog.Infof("Lookup revolved %s into %s", ip, resolvedIP)
	}
	return resolvedIP
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
			ApiPrefix:       GetIPUri(),
			IpEndPoints:     GetIpEndPoint(),
		}
	}

	return nfService
}

func GetIPUri() string {
	port := eirContext.SBIPort
	addr := eirContext.RegisterIP

	return fmt.Sprintf("%s://%s", eirContext.UriScheme, netip.AddrPortFrom(addr, uint16(port)).String())
}

func GetIpEndPoint() []models.IpEndPoint {
	if eirContext.RegisterIP.Is6() {
		return []models.IpEndPoint{
			{
				Ipv6Address: eirContext.RegisterIP.String(),
				Transport:   models.NrfNfManagementTransportProtocol_TCP,
				Port:        int32(eirContext.SBIPort),
			},
		}
	} else if eirContext.RegisterIP.Is4() {
		return []models.IpEndPoint{
			{
				Ipv4Address: eirContext.RegisterIP.String(),
				Transport:   models.NrfNfManagementTransportProtocol_TCP,
				Port:        int32(eirContext.SBIPort),
			},
		}
	}
	return nil
}

// Create new EIR context
func GetSelf() *EIRContext {
	return &eirContext
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
