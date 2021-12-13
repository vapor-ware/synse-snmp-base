package core

import (
	"testing"
	"time"

	"github.com/soniah/gosnmp"
	"github.com/stretchr/testify/assert"
)

func TestClient_Close(t *testing.T) {
	c := Client{
		GoSNMP: &gosnmp.GoSNMP{},
	}
	c.Close()
}

func TestNewClient(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(1024), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version3, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient2(t *testing.T) {
	// Same test, different configuration
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "tcp://localhost",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(161), client.Port)
	assert.Equal(t, "tcp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version2c, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient3(t *testing.T) {
	// Same test, different configuration
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v1",
		Agent:   "localhost:4321",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(4321), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version1, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient_NoSecurity(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	assert.Equal(t, "localhost", client.Target)
	assert.Equal(t, uint16(1024), client.Port)
	assert.Equal(t, "udp", client.Transport)
	assert.Equal(t, "", client.Community)
	assert.Equal(t, gosnmp.Version3, client.Version)
	assert.Equal(t, 1*time.Second, client.Timeout)
	assert.Equal(t, 1, client.Retries)
}

func TestNewClient_BadVersion(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "invalid-version",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadAgent(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "this is not a valid agent string!",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadTransport(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "something://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadPort(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:not-a-port",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_BadPort2(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:999999999999",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:   "NoAuthNoPriv",
			Context: "test",
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestNewClient_v2WithSecurity(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v2",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrNonV3SecurityParams, err)
}

func TestNewClient_BadSecLevel(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "invalid-level",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidMessageFlag, err)
}

func TestNewClient_BadAuthProtocol(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "invalid-protocol",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidAuthProtocol, err)
}

func TestNewClient_BadPrivProtocol(t *testing.T) {
	cfg := &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://localhost:1024",
		Timeout: 1 * time.Second,
		Retries: 1,
		Security: &SnmpV3Security{
			Level:    "AuthPriv",
			Context:  "test",
			Username: "test",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "test",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "invalid-protocol",
				Passphrase: "test",
			},
		},
	}

	client, err := NewClient(cfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, ErrInvalidPrivProtocol, err)
}

//
// Integration tests
//
// The integration tests run against an instance of the vaporio/snmp-emulator container
// configured with the UPS MIB (https://tools.ietf.org/html/rfc1628)
//

func getEmulatorClientConfig() *SnmpTargetConfiguration {
	return &SnmpTargetConfiguration{
		MIB:     "test-mib",
		Version: "v3",
		Agent:   "udp://127.0.0.1:1024",
		Security: &SnmpV3Security{
			Level:    "authPriv",
			Context:  "public",
			Username: "simulator",
			Authentication: &SnmpV3SecurityAuthentication{
				Protocol:   "SHA",
				Passphrase: "auctoritas",
			},
			Privacy: &SnmpV3SecurityPrivacy{
				Protocol:   "AES",
				Passphrase: "privatus",
			},
		},
	}
}

func TestClientGetExistingOidIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// OID of the UPS MIB upsIdent manufacturer
	oid := ".1.3.6.1.2.1.33.1.1.1.0"

	pdu, err := client.GetOid(oid)
	assert.NoError(t, err)

	assert.Equal(t, oid, pdu.Name)
	assert.Equal(t, gosnmp.OctetString, pdu.Type)
	assert.Equal(t, "Eaton Corporation", string(pdu.Value.([]uint8)))
}

func TestClientGetNonexistentOidIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// Some unknown OID
	oid := ".1.3.6.1.2.1.33.1.1.100.100"

	pdu, err := client.GetOid(oid)
	assert.NoError(t, err)

	assert.Equal(t, ".1.3.6.1.2.1.33.1.1.100.100", pdu.Name)
	assert.Equal(t, gosnmp.NoSuchInstance, pdu.Type)
	assert.Equal(t, nil, pdu.Value)
}

func TestClientGetBadOidIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// Some invalid OID
	oid := "."

	pdu, err := client.GetOid(oid)
	assert.EqualError(t, err, `marshal: unable to parse OID: strconv.Atoi: parsing "": invalid syntax`)
	assert.Nil(t, pdu)
}

func TestClientConnectToBadAgentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	cfg := getEmulatorClientConfig()
	cfg.Agent = "udp://1.2.3.4:9999"
	cfg.Timeout = 10 * time.Millisecond

	client, err := NewClient(cfg)
	defer client.Close()
	assert.NoError(t, err)

	// OID of the UPS MIB upsIdent manufacturer
	oid := ".1.3.6.1.2.1.33.1.1.1.0"

	pdu, err := client.GetOid(oid)
	assert.EqualError(t, err, "Request timeout (after 3 retries)")
	assert.Nil(t, pdu)
}

func TestClientGetSupportedDevicesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// The GetSupportedDevices call requires the client to be connected first,
	// so establish the connection.
	err = client.Connect()
	assert.NoError(t, err)

	// Root OID for the UPS MIB, used by the emulator.
	devices, err := client.GetSupportedDevices("1.3.6.1.2.1.33")
	assert.NoError(t, err)
	assert.Len(t, devices, 61)

	expectedOIDs := []string{
		"1.3.6.1.2.1.33.1.1.1.0",
		"1.3.6.1.2.1.33.1.1.2.0",
		"1.3.6.1.2.1.33.1.1.3.0",
		"1.3.6.1.2.1.33.1.1.4.0",
		"1.3.6.1.2.1.33.1.1.5.0",
		"1.3.6.1.2.1.33.1.1.6.0",
		"1.3.6.1.2.1.33.1.2.2.0",
		"1.3.6.1.2.1.33.1.2.3.0",
		"1.3.6.1.2.1.33.1.2.4.0",
		"1.3.6.1.2.1.33.1.2.5.0",
		"1.3.6.1.2.1.33.1.2.6.0",
		"1.3.6.1.2.1.33.1.2.7.0",
		"1.3.6.1.2.1.33.1.3.2.0",
		"1.3.6.1.2.1.33.1.3.3.1.2.1",
		"1.3.6.1.2.1.33.1.3.3.1.2.2",
		"1.3.6.1.2.1.33.1.3.3.1.2.3",
		"1.3.6.1.2.1.33.1.3.3.1.3.1",
		"1.3.6.1.2.1.33.1.3.3.1.3.2",
		"1.3.6.1.2.1.33.1.3.3.1.3.3",
		"1.3.6.1.2.1.33.1.3.3.1.4.1",
		"1.3.6.1.2.1.33.1.3.3.1.4.2",
		"1.3.6.1.2.1.33.1.3.3.1.4.3",
		"1.3.6.1.2.1.33.1.3.3.1.5.1",
		"1.3.6.1.2.1.33.1.3.3.1.5.2",
		"1.3.6.1.2.1.33.1.3.3.1.5.3",
		"1.3.6.1.2.1.33.1.4.1.0",
		"1.3.6.1.2.1.33.1.4.2.0",
		"1.3.6.1.2.1.33.1.4.3.0",
		"1.3.6.1.2.1.33.1.4.4.1.2.1",
		"1.3.6.1.2.1.33.1.4.4.1.2.2",
		"1.3.6.1.2.1.33.1.4.4.1.2.3",
		"1.3.6.1.2.1.33.1.4.4.1.3.1",
		"1.3.6.1.2.1.33.1.4.4.1.3.2",
		"1.3.6.1.2.1.33.1.4.4.1.3.3",
		"1.3.6.1.2.1.33.1.4.4.1.4.1",
		"1.3.6.1.2.1.33.1.4.4.1.4.2",
		"1.3.6.1.2.1.33.1.4.4.1.4.3",
		"1.3.6.1.2.1.33.1.4.4.1.5.1",
		"1.3.6.1.2.1.33.1.4.4.1.5.2",
		"1.3.6.1.2.1.33.1.4.4.1.5.3",
		"1.3.6.1.2.1.33.1.5.1.0",
		"1.3.6.1.2.1.33.1.5.2.0",
		"1.3.6.1.2.1.33.1.5.3.1.2.1",
		"1.3.6.1.2.1.33.1.5.3.1.2.2",
		"1.3.6.1.2.1.33.1.5.3.1.2.3",
		"1.3.6.1.2.1.33.1.6.1.0",
		"1.3.6.1.2.1.33.1.6.2.1.2.1",
		"1.3.6.1.2.1.33.1.6.2.1.2.2",
		"1.3.6.1.2.1.33.1.6.2.1.3.1",
		"1.3.6.1.2.1.33.1.6.2.1.3.2",
		"1.3.6.1.2.1.33.1.7.3.0",
		"1.3.6.1.2.1.33.1.7.4.0",
		"1.3.6.1.2.1.33.1.7.5.0",
		"1.3.6.1.2.1.33.1.7.6.0",
		"1.3.6.1.2.1.33.1.8.1.0",
		"1.3.6.1.2.1.33.1.8.5.0",
		"1.3.6.1.2.1.33.1.9.1.0",
		"1.3.6.1.2.1.33.1.9.2.0",
		"1.3.6.1.2.1.33.1.9.3.0",
		"1.3.6.1.2.1.33.1.9.4.0",
	}

	for i, oid := range expectedOIDs {
		assert.Contains(t, devices, oid, "oid:%s index:%d", oid, i)
		assert.Equal(t, struct{}{}, devices[oid])
	}
}

func TestClientGetSupportedDevicesNotConnectedIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// The GetSupportedDevices call requires the client to be connected first,
	// but this test checks the case where it is not connected, so do not
	// connect prior to calling GetSupportedDevices.

	// Root OID for the UPS MIB, used by the emulator.
	devices, err := client.GetSupportedDevices("1.3.6.1.2.1.33")
	assert.EqualError(t, err, "&GoSNMP.Conn is missing. Provide a connection or use Connect()")
	assert.Nil(t, devices)
}

func TestClientGetSupportedDevicesInvalidRootOIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test: --short flag set")
	}

	client, err := NewClient(getEmulatorClientConfig())
	defer client.Close()
	assert.NoError(t, err)

	// The GetSupportedDevices call requires the client to be connected first,
	// so establish the connection.
	err = client.Connect()
	assert.NoError(t, err)

	// Invalid root OID.
	devices, err := client.GetSupportedDevices("foo")
	assert.EqualError(t, err, `marshal: unable to parse OID: strconv.Atoi: parsing "foo": invalid syntax`)
	assert.Nil(t, devices)
}
