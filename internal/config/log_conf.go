package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var LogErr = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
	Level(zerolog.TraceLevel).
	With().
	Timestamp().
	Caller().
	Logger()

var LogInfo = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime}).
	Level(zerolog.TraceLevel).
	With().
	Timestamp().
	Logger()
