package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Ping is a handler for the ping endpoint
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
