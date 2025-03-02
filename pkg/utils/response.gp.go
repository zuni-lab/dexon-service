package utils

import (
	"github.com/labstack/echo/v4"
)

func OkResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, map[string]any{
		"ok":   true,
		"data": data,
	})
}
