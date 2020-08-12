package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-sdk/sdk"
	"github.com/vapor-ware/synse-snmp-base/pkg/core"
)

func TestReadHandlerFunc_NilDevice(t *testing.T) {
	readings, err := readHandlerFunc(nil)

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadHandlerFunc_NoAgent(t *testing.T) {
	readings, err := readHandlerFunc(&sdk.Device{
		Data: map[string]interface{}{},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadHandlerFunc_NoOid(t *testing.T) {
	readings, err := readHandlerFunc(&sdk.Device{
		Data: map[string]interface{}{
			"agent": "udp://localhost:1024",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadHandlerFunc_NoTargetCfg(t *testing.T) {
	readings, err := readHandlerFunc(&sdk.Device{
		Data: map[string]interface{}{
			"agent": "udp://localhost:1024",
			"oid":   "1.2.3.4",
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

func TestReadHandlerFunc_FailedClientInit(t *testing.T) {
	cfg := &core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "this:is-not@a:valid_agent+string",
	}

	readings, err := readHandlerFunc(&sdk.Device{
		Data: map[string]interface{}{
			"agent":      "udp://localhost:1024",
			"oid":        "1.2.3.4",
			"target_cfg": cfg,
		},
	})

	assert.Error(t, err)
	assert.Nil(t, readings)
}

//
// Integration tests
//
// The integration tests run against an instance of the vaporio/snmp-emulator container
// configured with the UPS MIB (https://tools.ietf.org/html/rfc1628)
//

func getEmulatorClientConfig() *core.SnmpTargetConfiguration {
	return &core.SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://127.0.0.1:1024",
		Security: &core.SnmpV3Security{
			Level:    "authPriv",
			Context:  "public",
			Username: "simulator",
			Authentication: &core.SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "auctoritas",
			},
			Privacy: &core.SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "privatus",
			},
		},
	}
}

func getReadTestDevice(oid string) *sdk.Device {
	cfg := getEmulatorClientConfig()
	return &sdk.Device{
		Type:   "test",
		Info:   "a test device",
		Output: "status", // some generic output for the test
		Data: map[string]interface{}{
			"mib":        "test-mib",
			"agent":      cfg.Agent,
			"target_cfg": cfg,
			"oid":        oid,
		},
		Context: map[string]string{
			"oid":   oid,
			"other": "foobar",
		},
	}
}

func TestReadHandlerFuncIdentStringIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	// OID of the UPS MIB upsIdent manufacturer
	oid := ".1.3.6.1.2.1.33.1.1.1.0"

	device := getReadTestDevice(oid)

	readings, err := readHandlerFunc(device)
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, "status", readings[0].Type) // determined by test device "output" type (status)
	assert.Equal(t, "Eaton Corporation", readings[0].Value)
	assert.Equal(t, map[string]string{
		"oid":   oid,
		"other": "foobar",
	}, readings[0].Context)
}

// FIXME (etd):  This is returning nil. Need to look into why that is. Is it not being parsed
//   correctly in the handler? Is it missing from the emulator? Something else?
// UPDATE (etd): Turning on logging, I see: got reading value for OID" name=.1.3.6.1.2.1.33.1.2.1.0 type=NoSuchInstance value="<nil>"
//   so I think this means 2 things: its not in the emulator, and we may need some sort of handling/logging for NoSuchInstance?
//func TestReadHandlerFuncBatteryEnumIntegration(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping integration test: --short flag set")
//	}
//
//	// OID of the UPS MIB upsBattery status
//	oid := ".1.3.6.1.2.1.33.1.2.1.0"
//
//	device := getReadTestDevice(oid)
//
//	readings, err := readHandlerFunc(device)
//	assert.NoError(t, err)
//	assert.Len(t, readings, 1)
//	assert.Equal(t, "status", readings[0].Type)  // determined by test device "output" type (status)
//	assert.Equal(t, "Eaton Corporation", readings[0].Value)
//	assert.Equal(t, map[string]string{
//		"oid": oid,
//		"other": "foobar",
//	}, readings[0].Context)
//}

func TestReadHandlerFuncBatteryTemperatureIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	// OID of the UPS MIB upsBattery temperature
	oid := ".1.3.6.1.2.1.33.1.2.7.0"

	device := getReadTestDevice(oid)

	readings, err := readHandlerFunc(device)
	assert.NoError(t, err)
	assert.Len(t, readings, 1)
	assert.Equal(t, "status", readings[0].Type) // determined by test device "output" type (status)
	assert.Equal(t, int(24), readings[0].Value)
	assert.Equal(t, map[string]string{
		"oid":   oid,
		"other": "foobar",
	}, readings[0].Context)
}
