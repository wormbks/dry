package cert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	realCertContent = `
-----BEGIN CERTIFICATE-----
MIIBtTCCAVqgAwIBAgIQV6ismg6yFAbnJm+JrnLC1zAKBggqhkjOPQQDAjAzMRMw
EQYDVQQKEwpxbG10ZWMuY29tMRwwGgYDVQQLDBNrMTJAREVTS1RPUC05N1YzNDRG
MB4XDTIzMDcxMDA5MDAwN1oXDTI0MDcwOTA5MDAwN1owMzETMBEGA1UEChMKcWxt
dGVjLmNvbTEcMBoGA1UECwwTazEyQERFU0tUT1AtOTdWMzQ0RjBZMBMGByqGSM49
AgEGCCqGSM49AwEHA0IABD572mT5PTvgZUwC9uLDvlMOD1ouv+obdDQ3KYGBdLON
I7C2koPcr2N7dJm88cTaXd6QUp5RTnVJ+4qkSs8snySjUDBOMA4GA1UdDwEB/wQE
AwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADAPBgNVHREECDAGhwR/AAABMAoGCCqGSM49BAMCA0kAMEYCIQDvp2ZmXOtI1MKN
PpL9n9NQ+JHm6DaHC29Sf/GRUja5pAIhAJZvpC+6Ja4fu2Ze+0Bujl9UjYAFl8A8
h4MkWTCDUcfs
-----END CERTIFICATE-----
`

	realKeyContent = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg1MFoOzit4PcSUrdv
3yRE4e8/Ey/iDy7k6G+9MpXTyZ6hRANCAAQ+e9pk+T074GVMAvbiw75TDg9aLr/q
G3Q0NymBgXSzjSOwtpKD3K9je3SZvPHE2l3ekFKeUU51SfuKpErPLJ8k
-----END PRIVATE KEY-----`
)

func Test_GetKeyPair(t *testing.T) {

	// Test case 1: real  Cert and key
	_, err := GetKeyPair(realCertContent, realKeyContent)
	assert.NoError(t, err)
}

func Test_CamControl_AddTlsCert_negative(t *testing.T) {

	// Mock certificate and key contents
	certContent := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
	keyContent := "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"

	// Mock file paths for certificate and key
	certPath := "/path/to/cert.pem"
	keyPath := "/path/to/key.pem"

	// Test case 1: Cert and key are file paths
	_, err := GetKeyPair(certPath, keyPath)
	assert.Error(t, err)

	// Test case 2: Cert is file path, key is content
	_, err = GetKeyPair(certPath, keyContent)
	assert.Error(t, err)

	// Test case 3: Cert is content, key is file path
	_, err = GetKeyPair(certContent, keyPath)
	assert.Error(t, err)

	// Test case 4: Cert and key are contents
	_, err = GetKeyPair(certContent, keyContent)
	assert.Error(t, err)

	// Test case 5: Error reading cert file
	_, err = GetKeyPair("/nonexistent/cert.pem", keyContent)
	assert.Error(t, err)

	// Test case 6: Error reading key file
	_, err = GetKeyPair(certContent, "/nonexistent/key.pem")
	assert.Error(t, err)

}

func Test_IsStringLikeFilePath(t *testing.T) {
	// Test case 1: Relative path
	result := IsStringLikeFilePath("../path/to/file.txt")
	assert.True(t, result)

	// Test case 2: Absolute path
	result = IsStringLikeFilePath("/path/to/file.txt")
	assert.True(t, result)

	// Test case 3: Path with redundant separators
	result = IsStringLikeFilePath("path/to//file.txt")
	assert.True(t, result)

	// Test case 4: Path with references to the current directory
	result = IsStringLikeFilePath("./path/to/file.txt")
	assert.True(t, result)

	// Test case 5: Path with directory separator
	result = IsStringLikeFilePath("path/to/file/")
	assert.True(t, result)

	// Test case 6: Path with only directory separator
	result = IsStringLikeFilePath("/")
	assert.True(t, result)

	// Test case 7: Empty string
	result = IsStringLikeFilePath("")
	assert.False(t, result)

	// Test case 8: Real  key
	result = IsStringLikeFilePath(realKeyContent)
	assert.False(t, result)
}

func TestSetTlsConfig_FilePaths(t *testing.T) {
	// Create a temporary certificate and key file
	certPath := "/tmp/cert.pem"
	keyPath := "/tmp/key.pem"
	certContent := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
	keyContent := "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
	err := os.WriteFile(certPath, []byte(certContent), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(keyPath, []byte(keyContent), 0644)
	assert.NoError(t, err)

	// Call the setTlsConfig method with file paths
	tlsConfig, err := GetTlsConfig(certPath, keyPath)

	// Assert that the TLS configuration is created correctly
	assert.NoError(t, err)
	assert.NotNil(t, tlsConfig)
	assert.Equal(t, certContent, string(tlsConfig.Certificates[0].Certificate[0]))
	assert.Equal(t, []byte(keyContent), (tlsConfig.Certificates[0].PrivateKey))

	// Cleanup the temporary files
	err = os.Remove(certPath)
	assert.NoError(t, err)
	err = os.Remove(keyPath)
	assert.NoError(t, err)
}

func TestSetTlsConfig_Strings(t *testing.T) {
	// Call the setTlsConfig method with certificate and key strings
	cert := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"
	key := "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"

	tlsConfig, err := GetTlsConfig(cert, key)

	// Assert that the TLS configuration is created correctly
	assert.NoError(t, err)
	assert.NotNil(t, tlsConfig)
	assert.Equal(t, cert, string(tlsConfig.Certificates[0].Certificate[0]))
	assert.Equal(t, []byte(key), tlsConfig.Certificates[0].PrivateKey)
}

func TestSetTlsConfig_InvalidFiles(t *testing.T) {
	// Call the setTlsConfig method with invalid file paths
	certPath := "/tmp/nonexistent_cert.pem"
	keyPath := "/tmp/nonexistent_key.pem"

	tlsConfig, err := GetTlsConfig(certPath, keyPath)

	// Assert that an error is returned and the TLS configuration is nil
	assert.Error(t, err)
	assert.Nil(t, tlsConfig)
}
