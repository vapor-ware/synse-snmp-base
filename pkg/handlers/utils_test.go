package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vapor-ware/synse-snmp-base/pkg/core"
)

func TestGetAgent(t *testing.T) {
	data := map[string]interface{}{
		"agent": "udp://localhost:1024",
	}

	agent, err := getAgent(data)
	assert.NoError(t, err)
	assert.Equal(t, "udp://localhost:1024", agent)
}

func TestGetAgent_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	agent, err := getAgent(data)
	assert.Error(t, err)
	assert.Equal(t, "", agent)
}

func TestGetAgent_BadType(t *testing.T) {
	data := map[string]interface{}{
		"agent": 1234,
	}

	agent, err := getAgent(data)
	assert.Error(t, err)
	assert.Equal(t, "", agent)
}

func TestGetOid(t *testing.T) {
	data := map[string]interface{}{
		"oid": "1.2.3.4",
	}

	oid, err := getOid(data)
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3.4", oid)
}

func TestGetOid_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	oid, err := getOid(data)
	assert.Error(t, err)
	assert.Equal(t, "", oid)
}

func TestGetOid_BadType(t *testing.T) {
	data := map[string]interface{}{
		"oid": 1234,
	}

	oid, err := getOid(data)
	assert.Error(t, err)
	assert.Equal(t, "", oid)
}

func TestGetTargetConfig(t *testing.T) {
	cfg := core.SnmpTargetConfiguration{}
	data := map[string]interface{}{
		"target_cfg": &cfg,
	}

	targetCfg, err := getTargetConfig(data)
	assert.NoError(t, err)
	assert.Equal(t, &cfg, targetCfg)
}

func TestGetTargetConfig_NotExist(t *testing.T) {
	data := map[string]interface{}{}

	targetCfg, err := getTargetConfig(data)
	assert.Error(t, err)
	assert.Nil(t, targetCfg)
}

func TestGetTargetConfig_BadType(t *testing.T) {
	data := map[string]interface{}{
		"target_cfg": 1234,
	}

	targetCfg, err := getTargetConfig(data)
	assert.Error(t, err)
	assert.Nil(t, targetCfg)
}

func TestParseEnum(t *testing.T) {
	data := map[string]interface{}{
		"enum": map[interface{}]interface{}{
			1: "foo",
			2: "bar",
		},
	}

	val, err := parseEnum(data, 1)
	assert.NoError(t, err)
	assert.Equal(t, "foo", val)
}

func TestParseEnum_NotAnEnum(t *testing.T) {
	data := map[string]interface{}{}

	val, err := parseEnum(data, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestParseEnum_BadEnumData(t *testing.T) {
	data := map[string]interface{}{
		"enum": "unexpected data",
	}

	val, err := parseEnum(data, 1)
	assert.Error(t, err)
	assert.Nil(t, val)
}

func TestParseEnum_NoEnumValue(t *testing.T) {
	data := map[string]interface{}{
		"enum": map[interface{}]interface{}{
			1: "foo",
			2: "bar",
		},
	}

	val, err := parseEnum(data, 3)
	assert.Error(t, err)
	assert.Nil(t, val)
}
