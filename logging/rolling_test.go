package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	// Define a sample configuration for testing.
	config := Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		Directory:             "/var/log",
		Filename:              "test.log",
		MaxSize:               10,
		MaxBackups:            3,
		MaxAge:                7,
		Compress:              false,
		LoggingLevel:          zerolog.DebugLevel,
	}

	// Call the Configure function.
	logger := Logger{} 
	logger.Configure(config)

	// Verify the logger is not nil.
	assert.NotNil(t, logger)
}

func TestSetLogLevel(t *testing.T) {
	// Define a log level for testing.
	level := zerolog.InfoLevel
	logger := Logger{} 

	// Call the SetLogLevel function.
	logger.setLogLevel(level)

	// Add assertions to check the zerolog settings as needed.
	// For example:
	assert.Equal(t, zerolog.TimeFormatUnixMs, zerolog.TimeFieldFormat)
	assert.Equal(t, level, zerolog.GlobalLevel())
}

func TestNewRollingFile(t *testing.T) {
	// Define a sample configuration for the rolling file.
	config := Config{
		FileLoggingEnabled: true,
		ConsoleLoggingEnabled: false,
		Directory:          "/var/log",
		Filename:           "test.log",
		MaxSize:            10,
		MaxBackups:         3,
		MaxAge:             7,
		Compress:           false,
		LoggingLevel: zerolog.WarnLevel,
	}

	logger := Logger{} 
	err := logger.Configure(config)
	assert.NoError(t, err)
	defer logger.Close()

	// Verify the closer is not nil.
	assert.NotNil(t, logger.logCloser)
	// Additional assertions can be added to further validate the writer's properties.
}


func Test_NewRollingFile_Error(t *testing.T) {
	// Define a sample configuration for the rolling file.
	config := Config{
		FileLoggingEnabled: true,
		ConsoleLoggingEnabled: false,
		Directory:          "/root/log",
		Filename:           "test.log",
		MaxSize:            10,
		MaxBackups:         3,
		MaxAge:             7,
		Compress:           false,
		LoggingLevel: zerolog.WarnLevel,
	}

	logger := Logger{} 
	err := logger.Configure(config)
	assert.Error(t, err)
	defer logger.Close()

	// Verify the closer is not nil.
	assert.Nil(t, logger.logCloser)
	// Additional assertions can be added to further validate the writer's properties.
}


func TestNewConsoleLoger(t *testing.T) {
	// Define a sample configuration for the rolling file.
	config := Config{
		FileLoggingEnabled: false,
		ConsoleLoggingEnabled: true,
		Directory:          "/var/log",
		Filename:           "test.log",
		MaxSize:            10,
		MaxBackups:         3,
		MaxAge:             7,
		Compress:           false,
		LoggingLevel: zerolog.WarnLevel,
	}

	logger := Logger{} 
	err := logger.Configure(config)
	assert.NoError(t, err)
	defer logger.Close()

	// Verify the closer is not nil.
	assert.Nil(t, logger.logCloser)
	// Additional assertions can be added to further validate the writer's properties.
	logger.Info().Msg("test")
}

func Test_NoLoggersEnabled(t *testing.T) {
	// Define a sample configuration for the rolling file.
	config := Config{
		FileLoggingEnabled: false,
		ConsoleLoggingEnabled: false,
		Directory:          "/var/log",
		Filename:           "test.log",
		MaxSize:            10,
		MaxBackups:         3,
		MaxAge:             7,
		Compress:           false,
	}

	logger := Logger{} 
	err := logger.Configure(config)
	assert.Error(t, err)
	defer logger.Close()

	// Verify the closer is not nil.
	assert.Nil(t, logger.logCloser)
	// Additional assertions can be added to further validate the writer's properties.
}
