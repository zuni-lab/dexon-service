package server

import (
	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/internal/chat"
	"github.com/zuni-lab/dexon-service/internal/health"
	"github.com/zuni-lab/dexon-service/internal/orders"
	"github.com/zuni-lab/dexon-service/internal/pools"
)

func setupRoute(e *echo.Echo) {
	api := e.Group("/api")
	health.Route(e, "/health")
	pools.Route(api, "/pools")
	chat.Route(api, "/chat")
	orders.Route(api, "/orders")
}
