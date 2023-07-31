package cert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetKeyPair(certPathOrCert string, keyPathOrKey string) (tls.Certificate, error) {
	fail := func(err error) (tls.Certificate, error) { return tls.Certificate{}, err }
	var cert []byte
	var key []byte
	var err error

	if IsStringLikeFilePath(certPathOrCert) {
		cert, err = os.ReadFile(certPathOrCert)
		if err != nil {
			err = fmt.Errorf("failed to load certificate file: %w", err)
			return fail(err)
		}
	} else {
		cert = []byte(certPathOrCert)
	}

	if IsStringLikeFilePath(keyPathOrKey) {
		key, err = os.ReadFile(keyPathOrKey)
		if err != nil {
			err = fmt.Errorf("failed to load key file: %w", err)
			return fail(err)
		}
	} else {
		key = []byte(keyPathOrKey)
	}

	certPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		err = fmt.Errorf("failed reading client certificate: %w", err)
		return fail(err)
	}
	// You add one or more certificates
	return certPair, err
}

func GetTlsConfig(certPathOrCert string, keyPathOrKey string) (tlsConfig *tls.Config, err error) {
	var cert []byte

	var key []byte

	if IsStringLikeFilePath(certPathOrCert) {
		cert, err = os.ReadFile(certPathOrCert)
		if err != nil {
			err = fmt.Errorf("failed to load certificate file: %s", err.Error())
			return nil, err
		}
	} else {
		cert = []byte(certPathOrCert)
	}

	if IsStringLikeFilePath(keyPathOrKey) {
		key, err = os.ReadFile(keyPathOrKey)
		if err != nil {
			err = fmt.Errorf("failed to load key file: %s", err.Error())
			return nil, err
		}
	} else {
		key = []byte(keyPathOrKey)
	}

	rootCA := x509.NewCertPool()
	rootCA.AppendCertsFromPEM(cert)

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  key,
			},
		},
		RootCAs: rootCA,
	}

	return tlsConfig, nil
}

// IsStringLikeFilePath checks if a string is similar to a file path.
//
// It takes a string as a parameter.
// It returns a boolean value.
func IsStringLikeFilePath(s string) bool {
	// Clean the path to remove any redundant
	// separators and references to the current directory
	cleanedPath := filepath.Clean(s)
	// Check if the cleaned path contains a directory separator
	containsSeparator := strings.ContainsAny(cleanedPath, string(filepath.Separator))
	containsBegin := strings.ContainsAny(cleanedPath, "BEGIN")
	// If the cleaned path is absolute or contains a directory separator, consider it as a file path
	return !containsBegin && (filepath.IsAbs(cleanedPath) || containsSeparator)
}
