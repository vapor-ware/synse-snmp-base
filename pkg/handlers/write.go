package handlers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
)

// writeHandlerFunc is the function which handles Writes for SNMP device handlers.
func writeHandlerFunc(device *sdk.Device, data *sdk.WriteData) error {
	// TODO: implement
	log.Error("[snmp] SNMP writes are not yet implemented")
	return fmt.Errorf("snmp writes not yet implmented in snmp plugin base")
}
