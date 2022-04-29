package proxy

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	UUID              uuid.UUID         `json:"uuid"`
	Method            string            `json:"method"`
	RequestURI        string            `json:"requestURI"`
	Headers           map[string]string `json:"headers"`
	CreatedAt         time.Time         `json:"createdAt"`
	Body              []byte            `json:"body"`
	BodyEscapedString string            `json:"bodyEscaped"`
	Response          *Response         `json:"response"`
}

type Response struct {
	Code              int               `json:"code"`
	Status            string            `json:"status"`
	Headers           map[string]string `json:"headers"`
	Body              []byte            `json:"body"`
	BodyEscapedString string            `json:"bodyEscaped"`
	Duration          time.Duration     `json:"duration"`
	DurationAsString  string            `json:"durationString"`
}

func FlattenHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)

	for key, values := range headers {
		result[key] = values[0]
	}

	return result
}
