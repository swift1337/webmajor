package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/swift1337/webmajor/internal/proxy"
	httpServer "github.com/swift1337/webmajor/internal/server/http"
	v1 "github.com/swift1337/webmajor/internal/server/http/v1"
	"github.com/swift1337/webmajor/internal/store"
	"github.com/swift1337/webmajor/web"
)

var (
	servicePort = flag.String("service-port", "8080", "incoming data port")
	proxyPort   = flag.String("proxy-port", "8000", "proxy destination port")
)

var logger zerolog.Logger

func main() {
	flag.Parse()
	logger = setupLogger()

	handler := v1.New(
		proxy.NewCaller("http://0.0.0.0:"+*proxyPort, logger),
		store.NewSyncSlice(),
		logger,
	)

	router := httpServer.NewRouter(
		httpServer.WithDashboardAPI("/__webmajor/api", handler),
		httpServer.WithDashboardAssets("/__webmajor", http.FS(web.Dashboard)),
		httpServer.WithProxy(handler),
	)

	logger.Info().Msg("starting http server")
	err := http.ListenAndServe(":"+*servicePort, router)

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
