package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/swift1337/webmajor/internal/log"
	"github.com/swift1337/webmajor/internal/proxy"
	httpServer "github.com/swift1337/webmajor/internal/server/http"
	v1 "github.com/swift1337/webmajor/internal/server/http/v1"
	"github.com/swift1337/webmajor/internal/store"
	"github.com/swift1337/webmajor/web"
)

const (
	ApiPath       = "/__webmajor/api"
	DashboardPath = "/__webmajor"
)

var (
	servicePort = flag.String("service-port", "8080", "WebMajor incoming port")
	sourceBase  = flag.String("source", "http://0.0.0.0:8000", "Destination service")
)

func main() {
	flag.Parse()
	logger := log.New()

	handler := v1.New(
		proxy.NewCaller(*sourceBase, time.Second*30, logger),
		store.NewSyncSlice(),
		logger,
	)

	router := httpServer.NewRouter(
		httpServer.WithDashboardAPI(ApiPath, handler),
		httpServer.WithDashboardAssets(DashboardPath, http.FS(web.DashboardFiles())),
		httpServer.WithProxy(handler),
	)

	logger.Info().Msg("starting server")
	logger.Info().Msgf("source base: %s", *sourceBase)
	logger.Info().Msgf("to visit dashboard, open http://localhost:%s%s", *servicePort, DashboardPath)

	err := http.ListenAndServe(":"+*servicePort, router)

	if err != nil {
		logger.Fatal().Err(err).Msg("error while running server")
	}
}
