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
