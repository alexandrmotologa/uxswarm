package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// SSEMessage represents an individual event sent over Server-Sent Events.
type SSEMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// SSEHub manages active client subscriptions for real-time test run streaming.
type SSEHub struct {
	mu          sync.RWMutex
	subscribers map[string]map[chan []byte]struct{}
}

// NewSSEHub initializes an event hub.
func NewSSEHub() *SSEHub {
	return &SSEHub{
		subscribers: make(map[string]map[chan []byte]struct{}),
	}
}

// Subscribe registers a client channel to receive events for a specific runID.
func (h *SSEHub) Subscribe(runID string) chan []byte {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan []byte, 32)
	if _, ok := h.subscribers[runID]; !ok {
		h.subscribers[runID] = make(map[chan []byte]struct{})
	}
	h.subscribers[runID][ch] = struct{}{}
	return ch
}

// Unsubscribe removes a client channel.
func (h *SSEHub) Unsubscribe(runID string, ch chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.subscribers[runID]; ok {
		delete(clients, ch)
		close(ch)
		if len(clients) == 0 {
			delete(h.subscribers, runID)
		}
	}
}

// Broadcast serializes and delivers an event to all subscribers of a run.
func (h *SSEHub) Broadcast(runID string, eventName string, data any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.subscribers[runID]
	if !ok || len(clients) == 0 {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", eventName, string(payload))
	msgBytes := []byte(msg)

	for ch := range clients {
		select {
		case ch <- msgBytes:
		default:
			// Non-blocking drop if client is saturated
		}
	}
}

// ServeHTTP handles incoming SSE connections at /api/runs/{id}/stream.
func (h *SSEHub) ServeHTTP(w http.ResponseWriter, r *http.Request, runID string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := h.Subscribe(runID)
	defer h.Unsubscribe(runID, ch)

	// Send initial connection ping
	fmt.Fprintf(w, "event: connected\ndata: {\"run_id\": %q}\n\n", runID)
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			_, _ = w.Write(msg)
			flusher.Flush()
		}
	}
}
