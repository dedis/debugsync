// Package sync defines the logger.
//
// sync is using a global logger with some default parameters.
// It is disabled by default and the level can be increased using
// an environment variable:
//
//	 SYNCLOG=trace
//	 SYNCLOG=info
//
// sync main feature is disabled by default and thus works seemingly
// like the original sync package from the standard. To enable the debugging
// feature, use the following environment variable, e.g:
//   SYNCON=true
//

package sync

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// EnvLogLevel is the name of the environment variable to change the logging
// level.
const EnvLogLevel = "DBGSYNCLOG"

// EnvDebugSwitch is the name of the environment variable to allow debugging.
const EnvDebugSwitch = "DBGSYNCON"

const defaultLevel = zerolog.NoLevel

func init() {
	dbg := os.Getenv(EnvDebugSwitch)
	DebugIsOn = strings.ToLower(dbg) == "true"

	lvl := os.Getenv(EnvLogLevel)

	var level zerolog.Level

	switch lvl {
	case "error":
		level = zerolog.ErrorLevel
	case "warn":
		level = zerolog.WarnLevel
	case "info":
		level = zerolog.InfoLevel
	case "debug":
		level = zerolog.DebugLevel
	case "trace":
		level = zerolog.TraceLevel
	case "":
		level = defaultLevel
	default:
		level = zerolog.Disabled
	}

	Logger = Logger.Level(level)
}

var logout = zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.RFC3339,
}

// Logger is a globally available logger instance. By default, it only prints
// error level messages but it can be changed through a environment variable.
var Logger = zerolog.New(logout).Level(defaultLevel).
	With().Timestamp().Logger().
	With().Caller().Logger()

// DebugIsOn allows to turn the debugging tool on
var DebugIsOn = false
