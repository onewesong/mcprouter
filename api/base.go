package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SSEWriter struct {
	w http.ResponseWriter
	f http.Flusher
}

func getSSEWriter(c echo.Context) (*SSEWriter, error) {
	w := c.Response().Writer
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("not support stream: %v", w)
	}

	return &SSEWriter{
		w: w,
		f: flusher,
	}, nil
}

func (s *SSEWriter) sendEventData(event string, data string) error {
	if _, err := fmt.Fprintf(s.w, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}

	s.f.Flush()

	return nil
}

func (s *SSEWriter) sendData(data string) error {
	if _, err := fmt.Fprintf(s.w, "data: %s\n\n", data); err != nil {
		return err
	}

	s.f.Flush()

	return nil
}
