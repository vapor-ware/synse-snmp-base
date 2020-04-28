package handlers

import "github.com/vapor-ware/synse-sdk/sdk"

// WriteOnly is an SNMP device handler for OIDs which are write-only.
var WriteOnly = sdk.DeviceHandler{
	Name: "write-only",
	Write: func(device *sdk.Device, data *sdk.WriteData) error {
		return nil
	},
}
