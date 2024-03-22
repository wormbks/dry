# Logging Package

This package provides a flexible logging solution with support for console and
file logging. It is built on top of `zerolog`, a fast and flexible
logging library for Go.

## Installation

To use this package, you can simply import it into your Go project:

```shell
go get -u github.com/bksworm/dry/logging
```

## Usage

```go
import (
    "github.com/bksworm/dry/logging"
)

func main() {
    // Configure logging
    logger := logging.Logger{}
    err := logger.Configure(logging.Config{
        ConsoleLoggingEnabled: true,
        FileLoggingEnabled:    true,
        Directory:             "./logs",
        Filename:              "app.log",
        MaxSize:               10,   // MB
        MaxBackups:            5,    // files
        MaxAge:                30,   // days
        Compress:              true,
        LoggingLevel:          zerolog.DebugLevel,
    })
    if err != nil {
        // Handle error
    }

    // Log messages
    logger.Info().Msg("This is an informational message")
    logger.Error().Msg("This is an error message")

    // Close the logger when done
    err = logger.Close()
    if err != nil {
        // Handle error
    }
}
```

## Configuration Options

- `ConsoleLoggingEnabled`: Enable console logging.
- `FileLoggingEnabled`: Enable file logging.
- `Directory`: Directory to log files.
- `Filename`: Name of the log file.
- `MaxSize`: Maximum size of each log file in megabytes before rolling.
- `MaxBackups`: Maximum number of rolled log files to keep.
- `MaxAge`: Maximum age of a log file in days before it is removed.
- `Compress`: Enable/disable log file compression.
- `LoggingLevel`: Set the logging level (`zerolog.Level`).

## Customization

You can customize the logging format and behavior by modifying the provided
configuration options and the `Configure` function. Additionally, you can extend
the functionality of this package according to your specific requirements.

## Credits

This package uses the following libraries:

- `zerolog`: Fast and flexible logging library for Go.
- `natefinch/lumberjack`: Rolling logger for Go.

## License

This package is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

For more information, visit the [GitHub repository](https://github.com/wormbks/dry/logging).
