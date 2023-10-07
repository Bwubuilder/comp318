package count

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type counter struct{}

type writeFlusher interface {
	http.ResponseWriter
	http.Flusher
}

func New() counter {
	return counter{}
}

func (c counter) send(wf writeFlusher, ctr int) {
	// Create event
	var evt bytes.Buffer
	evt.WriteString(fmt.Sprintf("event: count\ndata: %d\nid: %d\n\n", ctr, ctr))

	slog.Info("Sending", "msg", evt.String())

	// Send event
	wf.Write(evt.Bytes())
	wf.Flush()
}

func (c counter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check Last-Event-ID
	lastEventID := r.Header.Get("Last-Event-ID")
	slog.Info("ServeHTTP", "Last-Event-ID", lastEventID)

	// Convert ResponseWriter to a writeFlusher
	wf, ok := w.(writeFlusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	slog.Info("Converted to writeFlusher")

	// Set up event stream connection
	wf.Header().Set("Content-Type", "text/event-stream")
	wf.Header().Set("Cache-Control", "no-cache")
	wf.Header().Set("Connection", "keep-alive")
	wf.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Last-Event-ID")
	wf.Header().Set("Access-Control-Allow-Origin", "*")
	wf.WriteHeader(http.StatusOK)
	wf.Flush()

	slog.Info("Sent headers")

	ctr := 0

	if lastEventID != "" {
		val, err := strconv.Atoi(lastEventID)
		if err != nil {
			slog.Error("invalid last event id", "Last-Event-ID", lastEventID)
		} else {
			ctr = val
		}
	}

	// Send intial message
	c.send(wf, ctr)
	ctr++

	// Count forever
	for {
		select {
		case <-r.Context().Done():
			// Client closed connection
			slog.Info("Client closed connection")
			return
		case <-time.After(1 * time.Second):
			// Send next count
			c.send(wf, ctr)
			ctr++
		}
	}
}
