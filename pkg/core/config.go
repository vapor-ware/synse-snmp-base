package core

import (
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk/utils"
)

// SnmpTargetConfiguration defines the configuration for individual SNMP servers,
// which are defined in the plugin config's dynamicRegistration block.
type SnmpTargetConfiguration struct {
	MIB       string          `yaml:"mib,omitempty"`
	Version   string          `yaml:"version,omitempty"`
	Agent     string          `yaml:"agent,omitempty"`
	Community string          `yaml:"community,omitempty"`
	Timeout   time.Duration   `yaml:"timeout,omitempty"`
	Retries   int             `yaml:"retries,omitempty"`
	Security  *SnmpV3Security `yaml:"security,omitempty"`
}

// SnmpV3Security defines the security configuration for the SNMP connection. Only
// v3 of the SNMP protocol supports these security parameters.
type SnmpV3Security struct {
	Level          string                        `yaml:"level,omitempty"`
	Context        string                        `yaml:"context,omitempty"`
	Username       string                        `yaml:"username,omitempty"`
	Authentication *SnmpV3SecurityAuthentication `yaml:"authentication,omitempty"`
	Privacy        *SnmpV3SecurityPrivacy        `yaml:"privacy,omitempty"`
}

// SnmpV3SecurityAuthentication defines the authentication parameters for SNMP v3
// security, when auth is enabled.
type SnmpV3SecurityAuthentication struct {
	Protocol   string `yaml:"protocol,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
}

// SnmpV3SecurityPrivacy defines the privacy parameters for SNMP v3 security, when
// privacy is enabled.
type SnmpV3SecurityPrivacy struct {
	Protocol   string `yaml:"protocol,omitempty"`
	Passphrase string `yaml:"passphrase,omitempty"`
}

// LoadTargetConfiguration loads a map of data, which should be read in from the
// plugin config's dynamicRegistration block, into an SnmpTargetConfiguration struct.
func LoadTargetConfiguration(raw map[string]interface{}) (*SnmpTargetConfiguration, error) {
	var cfg SnmpTargetConfiguration

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &cfg,
	})
	if err != nil {
		log.Error("[snmp] failed to initialize config decoder")
		return nil, err
	}

	if err := decoder.Decode(raw); err != nil {
		log.WithFields(log.Fields{
			"data": utils.RedactPasswords(raw),
		}).Error("[snmp] failed decoding SNMP target configuration into struct")
		return nil, err
	}

	// Set defaults
	if cfg.Retries == 0 {
		cfg.Retries = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 3 * time.Second
	}

	return &cfg, nil
}
