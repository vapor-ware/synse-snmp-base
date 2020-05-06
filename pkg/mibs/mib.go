package mibs

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-base/pkg/core"
)

// MIB is a logical grouping of SnmpDevices which a SNMP plugin implementation
// should define. A single plugin implementation may define multiple MIBs. These
// MIBs are registered with the SNMP base plugin.
type MIB struct {
	Name    string
	RootOid string
	Devices []*SnmpDevice
}

// NewMIB creates a new MIB with the specified devices.
func NewMIB(name string, rootOid string, devices ...*SnmpDevice) *MIB {
	return &MIB{
		Name:    name,
		RootOid: rootOid,
		Devices: devices,
	}
}

// String returns a human-readable string, useful for identifying the
// MIB in logs.
func (mib *MIB) String() string {
	return fmt.Sprintf("[MIB %s (%s)]", mib.Name, mib.RootOid)
}

// LoadDevices loads Synse devices from the SNMP devices defined in the MIB.
func (mib *MIB) LoadDevices(cfg *core.SnmpTargetConfiguration, supported map[string]struct{}) ([]*sdk.Device, error) {
	if cfg == nil {
		return nil, errors.New("cannot load devices with nil SNMP target config")
	}

	log.WithFields(log.Fields{
		"mib":     mib.Name,
		"devices": len(mib.Devices),
	}).Debug("[snmp] loading devices for MIB")

	var devices []*sdk.Device
	for _, d := range mib.Devices {

		if _, exists := supported[d.OID]; !exists {
			log.WithFields(log.Fields{
				"oid":   d.OID,
				"agent": cfg.Agent,
			}).Debug("[snmp] mib device not supported by agent; will not load")
			continue
		}

		device, err := d.ToDevice()
		if err != nil {
			return nil, err
		}

		// Augment the device data with the MIB name and the target agent.
		// These pieces of information, along with the device OID (set in
		// the ToDevice call), are required by the plugin to generate a
		// unique ID for the device.
		device.Data["mib"] = mib.Name
		device.Data["agent"] = cfg.Agent
		device.Data["target_cfg"] = cfg

		devices = append(devices, device)
	}

	log.WithFields(log.Fields{"devices": devices}).Debug("[snmp] loaded devices")
	return devices, nil
}
