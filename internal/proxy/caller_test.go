package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

//nolint:errcheck
func TestCaller_Call(t *testing.T) {
	t.Run("InvalidURL", func(t *testing.T) {
		server := mockServer(func(w http.ResponseWriter, r *http.Request) {})
		caller := NewCaller("bla-bla-bla"+server.URL, time.Second, zerolog.Nop())
		incomingReq, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))

		proxyReq, err := caller.Call(incomingReq)

		assert.Error(t, err)
		assert.Nil(t, proxyReq)
	})

	t.Run("RequestTimeout", func(t *testing.T) {
		server := mockServer(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * 200)
			w.Write([]byte("hello"))
		})

		caller := NewCaller(server.URL, time.Millisecond*100, zerolog.Nop())
		incomingReq, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))

		proxyReq, err := caller.Call(incomingReq)

		assert.Error(t, err)
		assert.Nil(t, proxyReq)
	})

	t.Run("BrokenConnection", func(t *testing.T) {
		server := mockServer(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * 200)
			w.Write([]byte("hello"))
		})

		go func() {
			time.Sleep(time.Millisecond * 10)
			server.CloseClientConnections()
		}()

		c := NewCaller(server.URL, time.Millisecond*1000, zerolog.Nop())
		incomingReq, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))

		proxyReq, err := c.Call(incomingReq)

		assert.Error(t, err)
		assert.Nil(t, proxyReq)
	})

	t.Run("RequestCompleted", func(t *testing.T) {
		server := mockServer(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Hello", "world")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"user": "test"}`))
		})

		caller := NewCaller(server.URL, time.Millisecond*100, zerolog.Nop())
		incomingReq, _ := http.NewRequest("GET", "/", bytes.NewReader([]byte{}))

		proxyReq, err := caller.Call(incomingReq)

		assert.NoError(t, err)
		assert.Equal(t, proxyReq.Response.Body, []byte(`{"user": "test"}`))
		assert.Equal(t, proxyReq.Response.Headers["X-Hello"], "world")
	})
}

func mockServer(callback http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(callback)
}
