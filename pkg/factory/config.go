/*
 * EIR Configuration Factory
 */

package factory

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/adjivas/eir/internal/logger"
	"github.com/asaskevich/govalidator"
)

const (
	EirDefaultTLSKeyLogPath    = "./log/eirsslkey.log"
	EirDefaultCertPemPath      = "./cert/eir.pem"
	EirDefaultPrivateKeyPath   = "./cert/eir.key"
	EirDefaultConfigPath       = "./config/eircfg.yaml"
	EirSbiDefaultIP            = "127.0.0.7"
	EirSbiDefaultPort          = 8000
	EirSbiDefaultScheme        = "https"
	EirDefaultNrfUri           = "https://127.0.0.10:8000"
	EirDrResUriPrefix          = "/n5g-eir-eic/v1"
	EirMetricsDefaultPort      = 9091
	EirMetricsDefaultScheme    = "https"
	EirMetricsDefaultNamespace = "free5gc"
)

type DbType string

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
	Logger        *Logger        `yaml:"logger" valid:"required"`
	mu            sync.RWMutex
}

func (c *Config) Validate() (bool, error) {
	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Info struct {
	Version     string `yaml:"version,omitempty" valid:"required,in(1.1.0)"`
	Description string `yaml:"description,omitempty" valid:"type(string),optional"`
}

type Logger struct {
	Enable       bool   `yaml:"enable" valid:"type(bool)"`
	Level        string `yaml:"level" valid:"required,in(trace|debug|info|warn|error|fatal|panic)"`
	ReportCaller bool   `yaml:"reportCaller" valid:"type(bool)"`
}

const (
	EIR_DEFAULT_IP       = "127.0.0.4"
	EIR_DEFAULT_PORT     = "8000"
	EIR_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	Sbi             *Sbi     `yaml:"sbi" valid:"required"`
	Metrics         *Metrics `yaml:"metrics,omitempty" valid:"optional"`
	DefaultStatus   string   `yaml:"defaultStatus" valid:"in(WHITELISTED|BLACKLISTED),optional"`
	DbConnectorType DbType   `yaml:"dbConnectorType" valid:"required,in(mongodb)"`
	Mongodb         *Mongodb `yaml:"mongodb" valid:"optional"`
	NrfUri          string   `yaml:"nrfUri" valid:"url,required"`
	NrfCertPem      string   `yaml:"nrfCertPem,omitempty" valid:"optional"`
}

func (c *Configuration) validate() (bool, error) {
	if sbi := c.Sbi; sbi != nil {
		return sbi.validate()
	}

	if c.Metrics != nil {
		if _, err := c.Metrics.validate(); err != nil {
			return false, err
		}
		if c.Sbi != nil && c.Metrics.Port == c.Sbi.Port && c.Sbi.BindingIP == c.Metrics.BindingIPv4 {
			var errs govalidator.Errors
			err := fmt.Errorf("sbi and metrics bindings IP: %s and port: %d cannot be the same, please provide at least another port for the metrics", c.Sbi.BindingIP, c.Sbi.Port)
			errs = append(errs, err)
			return false, error(errs)
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Metrics struct {
	Scheme      string `yaml:"scheme" valid:"in(http|https)"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty" valid:"required,host"` // IP used to run the server in the node.
	Port        int    `yaml:"port,omitempty" valid:"required,port"`
	Tls         *Tls   `yaml:"tls,omitempty" valid:"optional"`
	Namespace   string `yaml:"namespace" valid:"optional"`
}

// This function is the mirror of the SBI one, I decided not to factor the code as it could in the future diverge.
// And it will reduce the cognitive overload when reading the function by not hiding the logic elsewhere.
func (m *Metrics) validate() (bool, error) {
	var errs govalidator.Errors

	if tls := m.Tls; tls != nil {
		if _, err := tls.validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if _, err := govalidator.ValidateStruct(m); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return false, error(errs)
	}

	return true, nil
}

type Sbi struct {
	Scheme     string `yaml:"scheme" valid:"in(http|https),optional"`
	RegisterIP string `yaml:"registerIP,omitempty" valid:"host,optional"` // IP that is registered at NRF.
	BindingIP  string `yaml:"bindingIP,omitempty" valid:"host,optional"`  // IP used to run the server in the node.
	Port       int    `yaml:"port" valid:"port,required"`
	Tls        *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

func (s *Sbi) validate() (bool, error) {
	// Set a default Schme if the Configuration does not provides one
	if s.Scheme == "" {
		s.Scheme = EirSbiDefaultScheme
	}

	// Set a default BindingIP/RegisterIP if the Configuration does not provides them
	if s.BindingIP == "" && s.RegisterIP == "" {
		s.BindingIP = EirSbiDefaultIP
		s.RegisterIP = EirSbiDefaultIP
	} else {
		// Complete any missing BindingIP/RegisterIP from RegisterIP/BindingIP
		if s.BindingIP == "" {
			s.BindingIP = s.RegisterIP
		} else if s.RegisterIP == "" {
			s.RegisterIP = s.BindingIP
		}
	}

	// Set a default Port if the Configuration does not provides one
	if s.Port == 0 {
		s.Port = EirSbiDefaultPort
	}

	if tls := s.Tls; tls != nil {
		if result, err := tls.validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(s)
	return result, err
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

func (t *Tls) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(t)
	return result, err
}

type Mongodb struct {
	Name string `yaml:"name" valid:"type(string),required"`
	Url  string `yaml:"url" valid:"required,required"`
}

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) SetLogEnable(enable bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Enable: enable,
			Level:  "info",
		}
	} else {
		c.Logger.Enable = enable
	}
}

func (c *Config) SetLogLevel(level string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level: level,
		}
	} else {
		c.Logger.Level = level
	}
}

func (c *Config) SetLogReportCaller(reportCaller bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		c.Logger = &Logger{
			Level:        "info",
			ReportCaller: reportCaller,
		}
	} else {
		c.Logger.ReportCaller = reportCaller
	}
}

func (c *Config) GetLogEnable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.Enable
}

func (c *Config) GetLogLevel() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return "info"
	}
	return c.Logger.Level
}

func (c *Config) GetLogReportCaller() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.ReportCaller
}

func (c *Config) GetCertPemPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Configuration.Sbi.Tls.Pem
}

func (c *Config) GetCertKeyPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Configuration.Sbi.Tls.Key
}

func (c *Config) GetMetricsScheme() string {
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.Scheme != "" {
		return c.Configuration.Metrics.Scheme
	}
	return EirMetricsDefaultScheme
}

func (c *Config) GetMetricsPort() int {
	if c.Configuration != nil && c.Configuration.Metrics != nil && c.Configuration.Metrics.Port != 0 {
		return c.Configuration.Metrics.Port
	}
	return EirMetricsDefaultPort
}

func (c *Config) GetMetricsBindingIP() string {
	bindIP := "0.0.0.0"
	if c.Configuration == nil || c.Configuration.Metrics == nil {
		return bindIP
	}

	if c.Configuration.Metrics.BindingIPv4 != "" {
		if bindIP = os.Getenv(c.Configuration.Metrics.BindingIPv4); bindIP != "" {
			logger.CfgLog.Infof("Parsing ServerIP [%s] from ENV Variable", bindIP)
		} else {
			bindIP = c.Configuration.Metrics.BindingIPv4
		}
	}
	return bindIP
}

func (c *Config) GetMetricsBindingAddr() string {
	return c.GetMetricsBindingIP() + ":" + strconv.Itoa(c.GetMetricsPort())
}

// We can see if there is a benefit to factor this tls key/pem with the sbi ones
func (c *Config) GetMetricsCertPemPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Configuration.Metrics != nil && c.Configuration.Metrics.Tls != nil {
		return c.Configuration.Metrics.Tls.Pem
	}

	return ""
}

func (c *Config) GetMetricsCertKeyPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Configuration.Metrics != nil && c.Configuration.Metrics.Tls != nil {
		return c.Configuration.Metrics.Tls.Key
	}

	return ""
}

func (c *Config) GetMetricsNamespace() string {
	if c.Configuration.Metrics != nil && c.Configuration.Metrics.Namespace != "" {
		return c.Configuration.Metrics.Namespace
	}
	return EirMetricsDefaultNamespace
}
