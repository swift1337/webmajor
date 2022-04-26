package proxy

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Code       int               `json:"code"`
	Method     string            `json:"method"`
	Headers    map[string]string `json:"headers"`
	Body       []byte            `json:"body"`
	RequestURI string            `json:"requestURI"`
	CreatedAt  time.Time         `json:"createdAt"`
	Duration   time.Duration     `json:"duration"`
}

func NewResponse(
	res *http.Response,
	method string,
	requestURI string,
	duration time.Duration,
) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)

	for key, values := range res.Header {
		headers[key] = values[0]
	}

	return &Response{
		Code:       res.StatusCode,
		Method:     method,
		RequestURI: requestURI,
		Headers:    headers,
		Body:       body,
		Duration:   duration,
		CreatedAt:  time.Now(),
	}, nil
}
