package proxy

import (
	"net/http"
	"time"
)

type Request struct {
	Method     string            `json:"method"`
	RequestURI string            `json:"requestURI"`
	Headers    map[string]string `json:"headers"`
	CreatedAt  time.Time         `json:"createdAt"`
	Body       []byte            `json:"body"`
	Response   *Response         `json:"response"`
}

type Response struct {
	Code     int               `json:"code"`
	Headers  map[string]string `json:"headers"`
	Body     []byte            `json:"body"`
	Duration time.Duration     `json:"duration"`
}

func FlattenHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)

	for key, values := range headers {
		result[key] = values[0]
	}

	return result
}
