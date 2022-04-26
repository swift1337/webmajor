package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/swift1337/webmajor/internal/proxy"
	v1 "github.com/swift1337/webmajor/internal/server/http/v1"
	"github.com/swift1337/webmajor/internal/store"
)

var (
	servicePort = flag.String("service-port", "8080", "incoming data port")
	proxyPort   = flag.String("proxy-port", "8000", "proxy destination port")
)

var logger zerolog.Logger

func main() {
	flag.Parse()

	logger = setupLogger()
	proxyCaller := proxy.NewCaller("http://0.0.0.0:"+*proxyPort, logger)

	handler := v1.New(
		proxyCaller,
		store.NewSyncSlice(),
		logger,
	)

	logger.Info().Msg("starting http server")
	err := http.ListenAndServe(":"+*servicePort, handler)

	if err != nil {
		logger.Fatal().Err(err).Msg("error while running http server")
	}
}

func setupLogger() zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.StampMilli,
	}

	return zerolog.New(writer).With().Timestamp().Logger()
}
