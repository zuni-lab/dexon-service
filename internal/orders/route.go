package orders

import (
	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/internal/orders/handlers"
)

func Route(g *echo.Group, path string) {
	// TODO: add middleware here

	middleware := echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// TODO: add middleware here
			return next(c)
		}
	})

	ordersGroup := g.Group(path, middleware)

	ordersGroup.GET("", handlers.List)
	ordersGroup.POST("", handlers.Create)
	ordersGroup.GET("/:id", handlers.Get)
	ordersGroup.POST("/:id/cancel", handlers.Cancel)
}
