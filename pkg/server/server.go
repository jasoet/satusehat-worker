package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var requestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "custom_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	},
)

type Operation func(e *echo.Echo)
type Shutdown func(e *echo.Echo)

func Start(port int, operation Operation, shutdown Shutdown) {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Info().
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("request")
			} else {
				log.Error().Err(v.Error).Msg("request error")
			}

			return nil
		},
	}))

	e.GET("/metrics", echoprometheus.NewHandler())

	if err := prometheus.Register(requestCounter); err != nil {
		log.Fatal().Err(err).Msg("failed to register request counter")
	}

	e.GET("/", func(c echo.Context) error {
		requestCounter.Inc()
		return c.String(http.StatusOK, "Home")
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go operation(e)

	go func() {
		log.Info().Msgf("Starting server, on port %d", port)
		if err := e.Start(fmt.Sprintf(":%v", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msg("gracefully shutting down")
	shutdown(e)

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown server")
	}

}
