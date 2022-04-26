package proxy

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Caller struct {
	client   *http.Client
	logger   zerolog.Logger
	basePath string
}

func NewCaller(
	basePath string,
	logger zerolog.Logger,
) *Caller {
	log := logger.With().Str("channel", "proxy_caller").Logger()

	return &Caller{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		basePath: basePath,
		logger:   log,
	}
}

func (c *Caller) Call(r *http.Request) (*Response, error) {
	destination := c.basePath + r.RequestURI

	proxyReq, err := http.NewRequest(r.Method, destination, r.Body)

	if err != nil {
		c.logger.Err(err).Msg("unable to create proxy request")
		return nil, err
	}

	startedAt := time.Now()
	res, err := http.DefaultClient.Do(proxyReq)

	duration := time.Now().Sub(startedAt)

	if err != nil {
		c.logger.Err(err).Str("destination", destination).Msg("request error")
		return nil, err
	}

	defer res.Body.Close()

	return NewResponse(
		res,
		r.Method,
		r.RequestURI,
		duration,
	)
}
