package server

import (
	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/internal/health"
)

func setupRoute(e *echo.Echo) {
	api := e.Group("/api")
	health.Route(e, "/health")
	_ = api
}
