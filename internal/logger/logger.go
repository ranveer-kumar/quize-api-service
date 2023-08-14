package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var consoleWriter zerolog.ConsoleWriter

func InitLogger() {
	consoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return "[" + strings.ToUpper(i.(string)) + "]"
	}

	log.Logger = log.Output(consoleWriter).Level(zerolog.InfoLevel).With().CallerWithSkipFrameCount(3).Str("service", "quize-api-service").Logger()
	// log.Logger = log.Level(zerolog.InfoLevel).With().CallerWithSkipFrameCount(3).Str("service", "quize-api-service").Logger()
	// sublogger := log.With().
	//              Str("component", "foo").
	//              Logger()

}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

// TODO: need to enable this with extra field in the logging for debugging purposes
// func WithFields(fields map[string]interface{}) zerolog.Logger {
//     newLogger := logger
//     for key, value := range fields {
//         newLogger = newLogger.With().Interface(key, value).Logger()
//     }
//     return newLogger
// }
