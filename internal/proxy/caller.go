package proxy

import (
	"bytes"
	"io"
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

func (c *Caller) Call(r *http.Request) (*Request, error) {
	destination := c.basePath + r.RequestURI

	defer r.Body.Close()

	// 1. Collect request body
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Err(err).Msg("unable to read request body")
		return nil, err
	}

	// 2. Init proxy request
	proxyReq, err := http.NewRequest(r.Method, destination, bytes.NewReader(requestBody))
	if err != nil {
		c.logger.Err(err).Msg("unable to create proxy request")
		return nil, err
	}

	// 3. Perform request
	createdAt := time.Now()

	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		c.logger.Err(err).Str("destination", destination).Msg("request error")
		return nil, err
	}

	defer res.Body.Close()

	duration := time.Now().Sub(createdAt)

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Err(err).Msg("unable to read response body")
		return nil, err
	}

	// 4. Construct result
	resultedRequest := &Request{
		Method:     r.Method,
		RequestURI: r.RequestURI,
		Headers:    FlattenHeaders(r.Header),
		CreatedAt:  createdAt,
		Body:       requestBody,
		Response: &Response{
			Code:     res.StatusCode,
			Headers:  FlattenHeaders(res.Header),
			Body:     responseBody,
			Duration: duration,
		},
	}

	return resultedRequest, nil
}
