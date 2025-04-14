/*
 * EIR Configuration Factory
 */

package factory

import (
	"fmt"
	"sync"

	"github.com/adjivas/eir/internal/logger"
	"github.com/asaskevich/govalidator"
)

const (
	EirDefaultTLSKeyLogPath  = "./log/eirsslkey.log"
	EirDefaultCertPemPath    = "./cert/eir.pem"
	EirDefaultPrivateKeyPath = "./cert/eir.key"
	EirDefaultConfigPath     = "./config/eircfg.yaml"
	EirSbiDefaultIP          = "127.0.0.7"
	EirSbiDefaultPort        = 8000
	EirSbiDefaultScheme      = "https"
	EirDefaultNrfUri         = "https://127.0.0.10:8000"
	EirDrResUriPrefix        = "/n5g-eir-eic/v1"
)

type DbType string

type Config struct {
	Info          *Info          `yaml:"info" valid:"required"`
	Configuration *Configuration `yaml:"configuration" valid:"required"`
	Logger        *Logger        `yaml:"logger" valid:"required"`
	sync.RWMutex
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

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
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
		errs = append(errs, fmt.Errorf("Invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()

	if c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}

func (c *Config) SetLogEnable(enable bool) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

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
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

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
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

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
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.Enable
}

func (c *Config) GetLogLevel() string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return "info"
	}
	return c.Logger.Level
}

func (c *Config) GetLogReportCaller() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	if c.Logger == nil {
		logger.CfgLog.Warnf("Logger should not be nil")
		return false
	}
	return c.Logger.ReportCaller
}

func (c *Config) GetCertPemPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Pem
}

func (c *Config) GetCertKeyPath() string {
	c.RLock()
	defer c.RUnlock()
	return c.Configuration.Sbi.Tls.Key
}
