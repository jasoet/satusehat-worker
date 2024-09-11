package main

import (
	"fmt"
	"github.com/jasoet/fhir-worker/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"os"
	"time"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Int("pid", os.Getpid()).
		Timestamp().
		Caller().Logger()

	rootCmd := app.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
