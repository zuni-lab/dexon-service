package handlers

import (
	"errors"
	"github.com/zuni-lab/dexon-service/internal/orders/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

func Cancel(c echo.Context) error {
	var (
		body services.CancelOrderBody
		err  error
		ctx  = c.Request().Context()
	)

	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	body.ID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid id"))
	}

	order, err := services.CancelOrder(ctx, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, order)
}
