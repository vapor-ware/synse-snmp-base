package core

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
)

// Errors relating to SNMP plugin client creation and usage.
var (
	ErrNonV3SecurityParams = errors.New("cannot define security parameters for SNMP versions other than v3")
)

// Client is a wrapper around a GoSNMP struct which adds some utility
// functions around it. Notably, it enables lazy connecting to the client,
// so the SNMP agent does not need to be reachable at plugin startup.
type Client struct {
	*gosnmp.GoSNMP

	isConnected bool
}

// GetOid gets the value for a specified OID.
func (c *Client) GetOid(oid string) (*gosnmp.SnmpPDU, error) {
	if !c.isConnected {
		log.Debug("[snmp] client establishing connection with agent")
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	result, err := c.Get([]string{oid})
	if err != nil {
		log.WithError(err).Error("[snmp] client failed to get OID")
		return nil, err
	}

	// Since we are currently only reading one OID, the result value will be
	// the first and only returned variable in the response.
	data := result.Variables[0]

	return &data, nil
}

// GetSupportedDevices gets all the OIDs for devices found on the target. This may not
// always be the full set of devices that a MIB defines.
//
// This returns a map of OIDs to empty struct. This map should be used during device creation
// to filter the MIB to only register those devices that a target supports. It is returned
// as a map to make OID lookups easier than iterating over a slice. Presence in the map means
// the device is supported, absence means it is not.
func (c *Client) GetSupportedDevices(rootOid string) (map[string]struct{}, error) {
	log.WithFields(log.Fields{
		"rootOid": rootOid,
	}).Debug("[snmp] getting supported devices for root OIO")

	results, err := c.BulkWalkAll(rootOid)
	if err != nil {
		log.WithError(err).Error("[snmp] failed to bulk walk all")
		return nil, err
	}

	log.WithFields(log.Fields{
		"size": len(results),
	}).Debug("[snmp] got bulk walk results")

	oids := make(map[string]struct{})
	for _, r := range results {
		oid := r.Name
		if strings.HasPrefix(oid, ".") {
			oid = oid[1:]
		}

		log.WithFields(log.Fields{
			"name":  oid,
			"value": r.Value,
			"type":  r.Type,
		}).Debug("[snmp] collecting walk result")
		oids[oid] = struct{}{}
	}

	return oids, nil
}

// Close the client connection.
func (c *Client) Close() {
	if c.Conn != nil {
		_ = c.Conn.Close()
	}
}

// NewClient creates a new instance of an SNMP Client for the given SNMP target
// configuration. The SNMP target configuration is defined in the dynamic configuration
// block for the plugin.
func NewClient(cfg *SnmpTargetConfiguration) (*Client, error) {
	// Verify that the configured version is valid.
	version, err := GetSNMPVersion(cfg.Version)
	if err != nil {
		return nil, err
	}

	log.WithField("version", version).Debug("[snmp] creating new client")

	// If configured to use SNMP v3, verify the security parameters.
	var securityModel gosnmp.SnmpV3SecurityModel
	var msgFlags gosnmp.SnmpV3MsgFlags
	var securityParams *gosnmp.UsmSecurityParameters
	var contextName string

	if version == gosnmp.Version3 {
		// If there are no security parameters defined, assume NoAuthNoPriv.
		if cfg.Security == nil {
			log.Info("[snmp] no security parameters defined, falling back to NoAuthNoPriv")

		} else {
			contextName = cfg.Security.Context

			msgFlags, err = GetSecurityFlags(cfg.Security.Level)
			if err != nil {
				return nil, err
			}
			if msgFlags != gosnmp.NoAuthNoPriv {
				securityModel = gosnmp.UserSecurityModel
			}

			var (
				authPass  = ""
				authProto = gosnmp.NoAuth
				privPass  = ""
				privProto = gosnmp.NoPriv
			)
			if cfg.Security.Authentication != nil {
				log.Debug("[snmp] parsing client auth configuration")
				authPass = cfg.Security.Authentication.Passphrase
				authProto, err = GetAuthProtocol(cfg.Security.Authentication.Protocol)
				if err != nil {
					return nil, err
				}
			}
			if cfg.Security.Privacy != nil {
				log.Debug("[snmp] parsing client privacy configuration")
				privPass = cfg.Security.Privacy.Passphrase
				privProto, err = GetPrivProtocol(cfg.Security.Privacy.Protocol)
				if err != nil {
					return nil, err
				}
			}
			securityParams = &gosnmp.UsmSecurityParameters{
				UserName:                 cfg.Security.Username,
				AuthenticationPassphrase: authPass,
				AuthenticationProtocol:   authProto,
				PrivacyPassphrase:        privPass,
				PrivacyProtocol:          privProto,
			}
		}

	} else {
		// If not using SNMP v3, no security parameters should be defined. If they
		// are, it should be considered a misconfiguration.
		if cfg.Security != nil {
			log.WithFields(log.Fields{
				"version": cfg.Version,
			}).Error("[snmp] security parameters unsupported for configured SNMP version")
			return nil, ErrNonV3SecurityParams
		}
	}

	agent := cfg.Agent
	if !strings.Contains(agent, "://") {
		agent = "udp://" + agent
	}

	u, err := url.Parse(agent)
	if err != nil {
		return nil, err
	}

	var transport string
	switch u.Scheme {
	case "tcp":
		transport = "tcp"
	case "", "udp":
		transport = "udp"
	default:
		return nil, fmt.Errorf("unsupported transport scheme: %s", u.Scheme)
	}

	portStr := u.Port()
	if portStr == "" {
		portStr = "161"
	}
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"host":      u.Hostname(),
		"port":      port,
		"transport": transport,
	}).Debug("[snmp] parsed client agent config")

	// Use a default timeout of 2s if one is not already set.
	timeout := 2 * time.Second
	if cfg.Timeout != 0 {
		timeout = cfg.Timeout
	}

	// Use a default of 3 retries if not already configured.
	retries := 3
	if cfg.Retries != 0 {
		retries = cfg.Retries
	}

	c := &Client{
		GoSNMP: &gosnmp.GoSNMP{
			Version:            version,
			Target:             u.Hostname(),
			Port:               uint16(port),
			Transport:          transport,
			Timeout:            timeout,
			Retries:            retries,
			Community:          cfg.Community,
			MsgFlags:           msgFlags,
			SecurityModel:      securityModel,
			SecurityParameters: securityParams,
			ContextName:        contextName,
			ExponentialTimeout: true,
			MaxOids:            gosnmp.MaxOids,
		},
	}

	return c, nil
}
