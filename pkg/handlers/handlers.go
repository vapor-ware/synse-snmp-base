package handlers

import (
	"github.com/vapor-ware/synse-sdk/sdk"
)

// ReadOnly is an SNMP device handler for OIDs which are read-only.
var ReadOnly = sdk.DeviceHandler{
	Name: "read-only",
	Read: readHandlerFunc,
}

// ReadWrite is an SNMP device handler for OIDs which are read-write.
var ReadWrite = sdk.DeviceHandler{
	Name:  "read-write",
	Read:  readHandlerFunc,
	Write: writeHandlerFunc,
}
