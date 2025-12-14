package logging

import (
	"strings"
	"sync"

	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

var (
	instance     *LoggingService
	once         sync.Once
	loggingMutex sync.RWMutex
)

type LoggingService struct {
	level zerolog.Level
}

// GetInstance returns the singleton logging service instance
func GetInstance() *LoggingService {
	once.Do(func() {
		// Check localStorage for persisted logging level
		level := zerolog.Disabled // Default to disabled
		if app.Window().Truthy() {
			storedValue := app.Window().Get("localStorage").Call("getItem", "rage_logging_level")
			if storedValue.Truthy() {
				level = parseLogLevel(storedValue.String())
			}
		}

		instance = &LoggingService{
			level: level,
		}

		// Set the initial global level based on stored preference
		zerolog.SetGlobalLevel(level)
	})
	return instance
}

// GetLevel returns the current log level
func (s *LoggingService) GetLevel() zerolog.Level {
	loggingMutex.RLock()
	defer loggingMutex.RUnlock()
	return s.level
}

// IsEnabled returns whether logging is currently enabled (not disabled)
func (s *LoggingService) IsEnabled() bool {
	loggingMutex.RLock()
	defer loggingMutex.RUnlock()
	return s.level != zerolog.Disabled
}

// SetLevel sets the logging level
func (s *LoggingService) SetLevel(level zerolog.Level) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	s.level = level

	zerolog.SetGlobalLevel(level)

	// Persist to localStorage
	if app.Window().Truthy() {
		app.Window().Get("localStorage").Call("setItem", "rage_logging_level", level.String())
	}
}

// SetEnabled enables or disables logging (for backward compatibility)
func (s *LoggingService) SetEnabled(enabled bool) {
	if enabled {
		s.SetLevel(zerolog.DebugLevel)
	} else {
		s.SetLevel(zerolog.Disabled)
	}
}

// Log returns a logger that respects the enabled state
func (s *LoggingService) Log(ctx app.Context) *zerolog.Logger {
	log := zerolog.Ctx(ctx)
	if !s.IsEnabled() {
		// Return a disabled logger
		disabled := zerolog.Nop()
		return &disabled
	}
	return log
}

// parseLogLevel converts a string to a zerolog.Level
func parseLogLevel(levelStr string) zerolog.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled", "disable", "off":
		return zerolog.Disabled
	default:
		return zerolog.Disabled
	}
}

// SetLogLevel sets the logging level from a string
func SetLogLevel(levelStr string) {
	level := parseLogLevel(levelStr)
	GetInstance().SetLevel(level)
	app.Window().Get("console").Call("log", "[RAGE] Logging level set to: "+level.String())
}

// EnableLogging is the function exposed to the browser console (for backward compatibility)
func EnableLogging(enabled bool) {
	GetInstance().SetEnabled(enabled)
	if enabled {
		app.Window().Get("console").Call("log", "[RAGE] Logging enabled (debug level)")
	} else {
		app.Window().Get("console").Call("log", "[RAGE] Logging disabled")
	}
}
