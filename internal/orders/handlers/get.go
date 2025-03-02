package handlers

import (
	"errors"
	"github.com/zuni-lab/dexon-service/internal/orders/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zuni-lab/dexon-service/pkg/utils"
)

func Get(c echo.Context) error {
	var (
		query services.GetOrderByIDQuery
		err   error

		ctx = c.Request().Context()
	)
	if err := utils.BindAndValidate(c, &query); err != nil {
		return err
	}

	query.ID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid id"))
	}

	order, err := services.GetOrderByID(ctx, query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return utils.OkResponse(c, http.StatusOK, order)
}
