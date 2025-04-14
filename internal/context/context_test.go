package context

import (
	"errors"
	"net/netip"
	"os"
	"testing"

	"github.com/adjivas/eir/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/stretchr/testify/assert"
)

func createConfigFile(t *testing.T, postContent []byte) *os.File {
	content := []byte(`info:
  version: 1.1.0
  description: EIR initial local configuration

logger:
  enable: true
  level: info

configuration:
  dbConnectorType: mongodb 
  mongodb:
    name: free5gc
    url: mongodb://localhost:27017
  nrfUri: http://127.0.0.10:8000
  nrfCertPem: cert/nrf.pem`)

	configFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("can't create temp file: %+v", err)
	}

	if _, err = configFile.Write(content); err != nil {
		t.Errorf("can't write content of temp file: %+v", err)
	}
	if _, err = configFile.Write(postContent); err != nil {
		t.Errorf("can't write content of temp file: %+v", err)
	}
	if err = configFile.Close(); err != nil {
		t.Fatal(err)
	}
	return configFile
}

func TestInitWithConfigIPv6(t *testing.T) {
	postContent := []byte(`

  sbi:
    scheme: http
    registerIP: 2001:db8::1:0:0:19
    bindingIP: 2001:db8::1:0:0:19
    port: 8000`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, eirContext.SBIPort, 8000)
	assert.Equal(t, eirContext.RegisterIP.String(), "2001:db8::1:0:0:19")
	assert.Equal(t, eirContext.BindingIP.String(), "2001:db8::1:0:0:19")
	assert.Equal(t, eirContext.UriScheme, models.UriScheme("http"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigIPv4(t *testing.T) {
	postContent := []byte(`
  sbi:
    scheme: http
    registerIP: "127.0.0.13"
    bindingIP: "127.0.0.13"
    port: 8131`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, eirContext.SBIPort, 8131)
	assert.Equal(t, eirContext.RegisterIP.String(), "127.0.0.13")
	assert.Equal(t, eirContext.BindingIP.String(), "127.0.0.13")
	assert.Equal(t, eirContext.UriScheme, models.UriScheme("http"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigEmptySBI(t *testing.T) {
	postContent := []byte("")

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	_, err := factory.ReadConfig(configFile.Name())
	assert.Equal(t, err, errors.New("Config validate Error"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigMissingRegisterIP(t *testing.T) {
	postContent := []byte(`
  sbi:
    bindingIP: "2001:db8::1:0:0:130"`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, eirContext.SBIPort, 8000)
	assert.Equal(t, eirContext.BindingIP.String(), "2001:db8::1:0:0:130")
	assert.Equal(t, eirContext.RegisterIP.String(), "2001:db8::1:0:0:130")
	assert.Equal(t, eirContext.UriScheme, models.UriScheme("https"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigMissingBindingIP(t *testing.T) {
	postContent := []byte(`
  sbi:
    registerIP: "2001:db8::1:0:0:131"`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, eirContext.SBIPort, 8000)
	assert.Equal(t, eirContext.BindingIP.String(), "2001:db8::1:0:0:131")
	assert.Equal(t, eirContext.RegisterIP.String(), "2001:db8::1:0:0:131")
	assert.Equal(t, eirContext.UriScheme, models.UriScheme("https"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigIPv6FromEnv(t *testing.T) {
	postContent := []byte(`
  sbi:
    scheme: http
    registerIP: "MY_REGISTER_IP"
    bindingIP: "MY_BINDING_IP"
    port: 8313`)

	configFile := createConfigFile(t, postContent)

	if err := os.Setenv("MY_REGISTER_IP", "2001:db8::1:0:0:130"); err != nil {
		t.Errorf("Can't set MY_REGISTER_IP env")
	}
	if err := os.Setenv("MY_BINDING_IP", "2001:db8::1:0:0:130"); err != nil {
		t.Errorf("Can't set MY_BINDING_IP env")
	}

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, eirContext.SBIPort, 8313)
	assert.Equal(t, eirContext.RegisterIP.String(), "2001:db8::1:0:0:130")
	assert.Equal(t, eirContext.BindingIP.String(), "2001:db8::1:0:0:130")
	assert.Equal(t, eirContext.UriScheme, models.UriScheme("http"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestResolveIPLocalhost(t *testing.T) {
	expectedAddr, err := netip.ParseAddr("::1")
	if err != nil {
		t.Errorf("invalid expected IP: %+v", expectedAddr)
	}

	addr := resolveIP("localhost")
	if addr != expectedAddr {
		t.Errorf("invalid IP: %+v", addr)
	}
	assert.Equal(t, addr, expectedAddr)
}

func TestResolveIPv4(t *testing.T) {
	expectedAddr, err := netip.ParseAddr("127.0.0.1")
	if err != nil {
		t.Errorf("invalid expected IP: %+v", expectedAddr)
	}

	addr := resolveIP("127.0.0.1")
	if addr != expectedAddr {
		t.Errorf("invalid IP: %+v", addr)
	}
}

func TestResolveIPv6(t *testing.T) {
	expectedAddr, err := netip.ParseAddr("2001:db8::1:0:0:1")
	if err != nil {
		t.Errorf("invalid expected IP: %+v", expectedAddr)
	}

	addr := resolveIP("2001:db8::1:0:0:1")
	if addr != expectedAddr {
		t.Errorf("invalid IP: %+v", addr)
	}
}

func TestInitWithConfigDefaultStatusBlack(t *testing.T) {
	postContent := []byte(`
  defaultStatus: "BLACKLISTED"
  sbi:
    scheme: http
    registerIP: "127.0.0.13"
    bindingIP: "127.0.0.13"
    port: 8131`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, "BLACKLISTED", eirContext.DefaultStatus)

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfigDefaultStatusWhite(t *testing.T) {
	postContent := []byte(`
  defaultStatus: "WHITELISTED"
  sbi:
    scheme: http
    registerIP: "127.0.0.13"
    bindingIP: "127.0.0.13"
    port: 8131`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	cfg, err := factory.ReadConfig(configFile.Name())
	if err != nil {
		t.Errorf("invalid read config: %+v %+v", err, cfg)
	}
	factory.EirConfig = cfg

	Init()

	assert.Equal(t, "WHITELISTED", eirContext.DefaultStatus)

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestInitWithConfiDefaultStatusgWrong(t *testing.T) {
	postContent := []byte(`
  defaultStatus: "PINKLISTED"
  sbi:
    scheme: http
    registerIP: "127.0.0.13"
    bindingIP: "127.0.0.13"
    port: 8131`)

	configFile := createConfigFile(t, postContent)

	// Test the initialization with the config file
	_, err := factory.ReadConfig(configFile.Name())
	assert.Equal(t, err, errors.New("Config validate Error"))

	// Close the config file
	t.Cleanup(func() {
		if err = os.RemoveAll(configFile.Name()); err != nil {
			t.Fatal(err)
		}
	})
}
