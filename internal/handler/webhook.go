package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/apialerts/hooks/internal/db"
	"github.com/go-chi/chi/v5"
)

const maxBodySize = 256 * 1024 // 256KB

func (h *Handler) receiveWebhook(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	endpoint, err := h.db.GetEndpoint(ctx, id)
	if err != nil {
		http.Error(w, "endpoint not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize+1))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	if len(body) > maxBodySize {
		http.Error(w, "payload too large (max 256KB)", http.StatusRequestEntityTooLarge)
		return
	}

	headers, _ := json.Marshal(r.Header)

	path := chi.URLParam(r, "*")
	if path == "" || path == "/" {
		path = ""
	}

	respStatus := endpoint.ResponseStatus
	responseBody, err := h.db.GetResponseBody(ctx, endpoint.ID, endpoint.ResponseStatus, endpoint.ResponseDelay)
	if err != nil {
		responseBody = defaultBodyForStatus(endpoint.ResponseStatus, endpoint.ResponseDelay)
	}

	seq, err := h.db.TouchEndpoint(ctx, id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	req := &db.Request{
		EndpointID:     endpoint.ID,
		Seq:            seq,
		Method:         r.Method,
		Path:           path,
		QueryParams:    r.URL.RawQuery,
		Headers:        headers,
		Body:           string(body),
		BodySize:       len(body),
		ContentType:    r.Header.Get("Content-Type"),
		SourceIP:       r.RemoteAddr,
		ResponseStatus: &respStatus,
		ResponseBody:   responseBody,
	}

	if err := h.db.CreateRequest(ctx, req); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// Trim to 50 requests max
	if seq > 50 {
		h.db.TrimRequests(ctx, endpoint.ID, 50)
	}

	if endpoint.ResponseDelay > 0 {
		time.Sleep(time.Duration(endpoint.ResponseDelay) * time.Millisecond)
	}

	if responseBody != "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(respStatus)
	w.Write([]byte(responseBody + "\n"))
}
