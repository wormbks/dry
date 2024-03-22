package logging

import (
	"errors"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ErrNoLoggerOutput = errors.New("no logger output enabled")

// Configuration for logging
type Config struct {
	// Enable console logging
	ConsoleLoggingEnabled bool `toml:"console_logging_enabled"`
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `toml:"file_logging_enabled"`
	// Directory to log to when file logging is enabled
	Directory string `toml:"directory"`
	// Filename is the name of the log file which will be placed inside the directory
	Filename string `toml:"filename"`
	// MaxSize the max size in MB of the log file before it's rolled
	MaxSize int `toml:"max_size"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `toml:"max_backups"`
	// MaxAge the max age in days to keep a log file
	MaxAge   int  `toml:"max_age"`
	Compress bool `toml:"compress"`
	// LoggingLevel sets the logging level
	LoggingLevel zerolog.Level
}



// Logger represents the logger
type Logger struct {
	*zerolog.Logger
	logCloser io.Closer
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.logCloser!= nil {
		return l.logCloser.Close()
	}
	return nil
}


// Configure configures the logger based on the provided configuration.
// It initializes console and/or file logging based on the settings in the
// provided Config struct.
func (l *Logger) Configure(config Config) error {
	var writers []io.Writer
	
	if !config.FileLoggingEnabled && !config.ConsoleLoggingEnabled {
		return ErrNoLoggerOutput
	}

	if config.FileLoggingEnabled {
		w, err := l.newRollingFile(config)
		if err != nil {
			log.Error().Err(err).Msg("can't create log file")
			return err
		}

		writers = append(writers, w)
	}

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "20060102-150405"})
	}

	mw := io.MultiWriter(writers...)

	var logger zerolog.Logger
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.CallerFieldName = "c"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	l.setLogLevel(config.LoggingLevel)
	if config.LoggingLevel <= zerolog.DebugLevel {
		logger = zerolog.New(mw).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(mw).With().Timestamp().Logger()
	}

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("consoleLogging", config.ConsoleLoggingEnabled).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Bool("compress", config.Compress).
		Msg("logging configured")

	l.Logger = &logger
	return nil
}

func (l *Logger) setLogLevel(level zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	if level <= zerolog.DebugLevel {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
			return file + ":" + strconv.Itoa(line)
		}
	} else {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return ""
		}
	}

	zerolog.SetGlobalLevel(level)
}

func (l *Logger) newRollingFile(config Config) (io.Writer, error) {
	err := os.MkdirAll(config.Directory, 0o755) // #nosec G301
	if err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil, err
	}
	lj := &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		Compress:   config.Compress,
	}

	l.logCloser = lj

	wr := diode.NewWriter(lj, 1000, 100*time.Millisecond, func(missed int) {
		// NOTE: it is very hard to write overlading test  for zerolog logger.
		log.Printf("Dropped %d messages", missed)
	})

	return wr, err
}

