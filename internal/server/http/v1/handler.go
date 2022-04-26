package v1

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/swift1337/webmajor/internal/proxy"
	"github.com/swift1337/webmajor/internal/store"
)

type Handler struct {
	proxyCaller  *proxy.Caller
	requestStore *store.SyncSlice
	logger       zerolog.Logger
}

func New(
	proxyCaller *proxy.Caller,
	requestStore *store.SyncSlice,
	logger zerolog.Logger,
) *Handler {
	log := logger.With().Str("channel", "handler_v1").Logger()

	return &Handler{
		proxyCaller:  proxyCaller,
		requestStore: requestStore,
		logger:       log,
	}
}

// ServeHTTP by default with handler calls HandleProxyRequest
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.HandleProxyRequest(w, r)
}

func (h *Handler) HandleProxyRequest(w http.ResponseWriter, r *http.Request) {
	result, err := h.proxyCaller.Call(r)

	if err != nil {
		h.logger.Err(err).Msg("error while proxying request")
		writeErr(w)
		return
	}

	h.logger.Info().
		Str("destination", result.RequestURI).
		Int("statusCode", result.Code).
		Str("method", result.Method).
		Dur("duration", result.Duration).
		Int("body_size", len(result.Body)).
		Msg("got response from destination")

	writeResponse(w, result)

	h.requestStore.Append(result)
}

// ListRequests get all requests that are made
func (h *Handler) ListRequests(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	result := make([]*proxy.Response, 0)

	for item := range h.requestStore.Iterate() {
		responseItem, ok := item.Value.(*proxy.Response)

		if !ok {
			h.logger.Warn().Msg("unable to cast slice value to responseItem")
			continue
		}

		result = append(result, responseItem)
	}

	encoded, err := json.Marshal(result)

	if err != nil {
		h.logger.Err(err).Msg("unable to marshal response")
		return
	}

	w.Write(encoded)
}

func writeResponse(w http.ResponseWriter, proxyResponse *proxy.Response) {
	for key, value := range proxyResponse.Headers {
		w.Header().Add(key, value)
	}

	w.WriteHeader(proxyResponse.Code)
	w.Write(proxyResponse.Body)
}

func writeErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("error while requesting destination"))
}
