package sec

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	once       sync.Once
	sharedCert *tls.Certificate
	sharedRoot *x509.CertPool
	sharedTLS  *tls.Config
)

var (
	ErrSTlsConfigNotInitialized = errors.New("tls config not initialized")
	ErrKeysEmpy                 = errors.New("both certPath and keyPath must be provided or left empty")
	ErrEmptyPaths               = errors.New("all tls paths are empty")
)

// InitTLSConfig initializes the shared TLS configuration.
// using the provided CA, certificate, and key paths. The once package from
// the sync package ensures that the initialization is performed only once
// even if multiple goroutines call it concurrently.
func InitTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	var errInit error

	once.Do(func() {
		sharedRoot, errInit = GetCA(caPath)
		if errInit != nil {
			return
		}

		certPair, errInit := GetKeyPair(certPath, keyPath)
		if errInit != nil {
			return
		}
		sharedCert = &certPair

		sharedTLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
			Certificates: []tls.Certificate{
				*sharedCert,
			},
			RootCAs: sharedRoot,
		}
	})

	return sharedTLS, errInit
}

// GetSharedTLSConfig returns the shared TLS configuration.
func GetSharedTLSConfig() (*tls.Config, error) {
	if sharedTLS == nil {
		return nil, ErrSTlsConfigNotInitialized
	}
	return sharedTLS, nil
}

func GetCA(caPathOrCert string) (*x509.CertPool, error) {
	// Load server CA certificate
	var cert []byte
	var err error
	if isStringLikeFilePath(caPathOrCert) {
		cert, err = os.ReadFile(filepath.Clean(caPathOrCert))
		if err != nil {
			return nil, fmt.Errorf("failed to load ca: %w", err)
		}
	} else {
		cert = []byte(caPathOrCert)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to append ca: %w", err)
	}

	return caCertPool, nil
}

// GetKeyPair returns a TLS certificate and an error.
//
// It takes two parameters:
//   - certPathOrCert: a string representing either the path to a certificate file or the certificate
//     content itself.
//   - keyPathOrKey: a string representing either the path to a key file or the key content itself.
//
// It returns a tls.Certificate and an error. The function attempts to load the certificate and key
// files if they are specified as file paths. If the files cannot be loaded, an error is returned.
// Otherwise, the function uses the provided certificate and key content to create a tls.Certificate
// and returns it along with any error that occurred during the process.
func GetKeyPair(certPathOrCert string, keyPathOrKey string) (tls.Certificate, error) {
	fail := func(err error) (tls.Certificate, error) { return tls.Certificate{}, err }
	var cert []byte
	var key []byte
	var err error

	if isStringLikeFilePath(certPathOrCert) {
		cert, err = os.ReadFile(filepath.Clean(certPathOrCert))
		if err != nil {
			err = fmt.Errorf("failed to load certificate file: %w", err)
			return fail(err)
		}
	} else {
		cert = []byte(certPathOrCert)
	}

	if isStringLikeFilePath(keyPathOrKey) {
		key, err = os.ReadFile(filepath.Clean(keyPathOrKey))
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

// GetTlsConfig returns a client TLS configuration based on the provided certificate and key.
//
// It takes two parameters:
//   - certPathOrCert: a string representing the path to the certificate file or the certificate itself.
//   - keyPathOrKey: a string representing the path to the key file or the key itself.
//
// It returns a *tls.Config object and an error.
func GetTlsConfig(certPathOrCert string, keyPathOrKey string) (tlsConfig *tls.Config, err error) {
	var cert []byte

	var key []byte

	if isStringLikeFilePath(certPathOrCert) {
		cert, err = os.ReadFile(filepath.Clean(certPathOrCert))
		if err != nil {
			err = fmt.Errorf("failed to load certificate file: %w", err)
			return nil, err
		}
	} else {
		cert = []byte(certPathOrCert)
	}

	if isStringLikeFilePath(keyPathOrKey) {
		key, err = os.ReadFile(filepath.Clean(keyPathOrKey))
		if err != nil {
			err = fmt.Errorf("failed to load key file: %ws", err)
			return nil, err
		}
	} else {
		key = []byte(keyPathOrKey)
	}

	rootCA := x509.NewCertPool()
	rootCA.AppendCertsFromPEM(cert)

	tlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
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

// ReadTLSConfig reads the TLS certificate and key files and returns a TLS config.
// If all file paths are empty, it returns nil with no error.
// If only one of certPath or keyPath is provided, it returns an error.
// Otherwise, it loads the cert/key pair and CA cert if provided,
// and returns a TLS config.
func ReadTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	// Check if all strings are empty
	if caPath == "" && certPath == "" && keyPath == "" {
		// If all paths are empty, return an empty TLS config
		return nil, ErrEmptyPaths
	}

	// Check if either certPath or keyPath is empty
	if (certPath == "" && keyPath != "") || (certPath != "" && keyPath == "") {
		// If one of certPath or keyPath is empty, return an error
		return nil, ErrKeysEmpy
	}

	// If caPath is not empty, call GetCA to get the CA certificate
	var sharedRoot *x509.CertPool
	var err error
	if caPath != "" {
		sharedRoot, err = GetCA(caPath)
		if err != nil {
			return nil, err
		}
	}

	// If both certPath and keyPath are not empty, call GetKeyPair to get the certificate pair
	var certPair tls.Certificate
	if certPath != "" && keyPath != "" {
		certPair, err = GetKeyPair(certPath, keyPath)
		if err != nil {
			return nil, err
		}
	}

	// Create the TLS config
	sharedTLS := &tls.Config{
		MinVersion: tls.VersionTLS12,
		Certificates: []tls.Certificate{
			certPair,
		},
		RootCAs: sharedRoot,
	}

	return sharedTLS, nil
}

// isStringLikeFilePath checks if a string is similar to a file path.
//
// It takes a string as a parameter.
// It returns a boolean value.
func isStringLikeFilePath(s string) bool {
	// Clean the path to remove any redundant
	// separators and references to the current directory
	cleanedPath := filepath.Clean(s)
	// Check if the cleaned path contains a directory separator
	containsSeparator := strings.ContainsAny(cleanedPath, string(filepath.Separator))
	containsBegin := strings.ContainsAny(cleanedPath, "BEGIN")
	// If the cleaned path is absolute or contains a directory separator, consider it as a file path
	return !containsBegin && (filepath.IsAbs(cleanedPath) || containsSeparator)
}
