package handlers

import (
	"github.com/zuni-lab/dexon-service/internal/orders/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

func List(c echo.Context) error {
	ctx := c.Request().Context()

	var query services.ListOrdersByWalletQuery
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}

	orders, err := services.ListOrderByWallet(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, orders)
}
