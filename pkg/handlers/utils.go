package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-snmp-base/pkg/core"
)

// getAgent is a convenience function to safely get the "agent" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "agent" field is expected to exist in the device Data and it is expected to
// be a string, this function returns an error if it does not exist or cannot be cast to
// a string.
func getAgent(data map[string]interface{}) (string, error) {
	agentIface, exists := data["agent"]
	if !exists {
		return "", fmt.Errorf("expected field 'agent' in device data, but not found")
	}
	agent, ok := agentIface.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast 'agent' value (%T) to string", agentIface)
	}
	return agent, nil
}

// getOid is a convenience function to safely get the "oid" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "oid" field is expected to exist in the device Data and it is expected to
// be a string, this function returns an error if it does not exist or cannot be cast to
// a string.
func getOid(data map[string]interface{}) (string, error) {
	oidIface, exists := data["oid"]
	if !exists {
		return "", fmt.Errorf("expected field 'oid' in device data, but not found")
	}
	oid, ok := oidIface.(string)
	if !ok {
		return "", fmt.Errorf("failed to cast 'oid' value (%T) to string", oidIface)
	}
	return oid, nil
}

// getTargetConfig is a convenience function to safely get the "target_cfg" value out of a device's
// Data field and cast it to the appropriate type.
//
// Since the "target_cfg" field is expected to exist in the device Data and it is expected to
// be an SnmpTargetConfiguration, this function returns an error if it does not exist or cannot
// be cast to an SnmpTargetConfiguration.
func getTargetConfig(data map[string]interface{}) (*core.SnmpTargetConfiguration, error) {
	cfgIface, exists := data["target_cfg"]
	if !exists {
		return nil, fmt.Errorf("expected field 'target_cfg' in device data, but not found")
	}
	cfg, ok := cfgIface.(*core.SnmpTargetConfiguration)
	if !ok {
		return nil, fmt.Errorf("failed to cast 'target_cfg' value (%T) to SnmpTargetConfiguration", cfgIface)
	}
	return cfg, nil
}

// parseEnum checks to see if the device value is an enumeration, and if so, converts
// the value to the corresponding enumeration value based on a lookup table defined
// in the device Data.
func parseEnum(data map[string]interface{}, value interface{}) (interface{}, error) {
	enumIface, isEnum := data["enum"]
	if !isEnum {
		return value, nil
	}

	log.WithFields(log.Fields{
		"value": value,
		"enum":  enumIface,
	}).Debug("[snmp] device value is an enumeration")
	enumMap, ok := enumIface.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("enumeration for device value is not properly defined (%T)", enumIface)
	}
	val, exists := enumMap[value]
	if !exists {
		log.WithFields(log.Fields{
			"map":   enumMap,
			"value": value,
		}).Error("[snmp] device enum value does not exist in lookup")
		return nil, fmt.Errorf("device value does not exist in enum map")
	}
	log.WithFields(log.Fields{
		"value": value,
	}).Debug("[snmp] using enumeration value")
	return val, nil
}
