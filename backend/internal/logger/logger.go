package logger

import (
	"github.com/rs/zerolog"
	"os"
)

var Log zerolog.Logger

func Init(env string) {

	if env == "dev" {
		Log = zerolog.New(
			zerolog.ConsoleWriter{
				Out: os.Stdout,
			},
		).With().
			Timestamp().
			Logger()

		return
	}

	Log = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}
