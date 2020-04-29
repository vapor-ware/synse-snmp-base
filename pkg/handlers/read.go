package handlers

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/soniah/gosnmp"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-sdk/sdk/output"
	"github.com/vapor-ware/synse-snmp-base/pkg/core"
)

// readHandlerFunc is the function which handles Reads for SNMP device handlers.
func readHandlerFunc(device *sdk.Device) ([]*output.Reading, error) {
	if device == nil {
		return nil, errors.New("unable to read from nil device")
	}

	// Get data cached in device.Data
	agent, err := getAgent(device.Data)
	if err != nil {
		return nil, err
	}
	oid, err := getOid(device.Data)
	if err != nil {
		return nil, err
	}
	targetConfig, err := getTargetConfig(device.Data)
	if err != nil {
		return nil, err
	}

	// Create a new client with the target configuration.
	c, err := core.NewClient(targetConfig)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	log.WithFields(log.Fields{
		"agent": agent,
		"oid":   oid,
	}).Debug("[snmp] reading OID")

	result, err := c.GetOid(oid)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"value": result.Value,
		"name":  result.Name,
		"type":  result.Type,
	}).Debug("[snmp] got reading value for OID")

	var value interface{}
	switch result.Type {
	case gosnmp.OctetString:
		ascii, err := core.BytesIfaceToASCII(result.Value)
		if err != nil {
			return nil, err
		}
		value = ascii
	default:
		value = result.Value
	}

	// Check if the device has enumerated values. If so, an "enum" map is present
	// in the device Data. This is set via the device config.
	value, err = parseEnum(device.Data, value)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"value": value,
	}).Debug("[snmp] final value")

	o := output.Get(device.Output)
	if o == nil {
		return nil, fmt.Errorf("unable to format reading: device output not defined")
	}

	return []*output.Reading{
		o.MakeReading(value),
	}, nil
}
