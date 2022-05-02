package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New() zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.StampMilli,
	}

	return zerolog.New(writer).With().Timestamp().Logger()
}
