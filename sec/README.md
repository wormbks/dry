# sec Package

The `sec` package provides utilities for managing TLS configurations and
certificates in Go. It includes functions for initializing TLS configurations,
loading certificate files, and creating TLS configurations for secure communication.

## Installation

To use this package, you can import it into your Go project:

```shell
go get -u github.com/wormbks/dry/sec
```

## Usage

```go
import (
    "github.com/yourusername/sec"
)

func main() {
    // Initialize TLS configuration
    tlsConfig, err := sec.InitTLSConfig("ca.crt", "cert.crt", "key.key")
    if err != nil {
        // Handle error
    }

    // Use tlsConfig for secure communication
}
```

## Functions

### InitTLSConfig

```go
func InitTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error)
```

`InitTLSConfig` initializes a TLS configuration using the provided CA, certificate,
and key paths. It returns a `*tls.Config` object and an error. This function
ensures that the initialization is performed only once, even if called
concurrently by multiple goroutines.

### GetSharedTLSConfig

```go
func GetSharedTLSConfig() (*tls.Config, error)
```

`GetSharedTLSConfig` returns the shared TLS configuration initialized by
`InitTLSConfig`. If the TLS configuration is not initialized, it returns an error.

### ReadTLSConfig

```go
func ReadTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error)
```

`ReadTLSConfig` reads the TLS certificate and key files and returns a TLS
configuration. If all file paths are empty, it returns `nil` with no error.
If only one of the certificate or key paths is provided, it returns an error.

### GetTlsConfig

```go
func GetTlsConfig(certPathOrCert string, keyPathOrKey string) (*tls.Config, error)
```

`GetTlsConfig` returns a client TLS configuration based on the provided certificate and key.

## Customization

You can customize the package according to your specific TLS configuration
requirements by using the provided functions and options.

## Error Handling

The package returns appropriate errors when encountering issues such as missing
paths, invalid certificates, or keys.

## Credits

This package utilizes the standard Go `crypto/tls` and `crypto/x509` libraries
for TLS communication and certificate management.

## License

This package is licensed under the MIT License. See the [LICENSE](./LICENSE)
file for details.

For more information, visit the [GitHub repository](https://github.com/wormbks/dry/sec).


