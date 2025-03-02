package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type ErrResponse struct {
	Ok       bool        `json:"ok"`
	Message  string      `json:"message"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func HttpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	if m, ok := err.(*ValidationError); ok {
		err := c.JSON(http.StatusBadRequest, m)
		if err != nil {
			return
		}
	} else if m, ok := err.(*echo.HTTPError); ok {
		switch mType := m.Message.(type) {
		case string:
			err := c.JSON(m.Code, ErrResponse{Message: mType})
			if err != nil {
				return
			}
		case error:
			err := c.JSON(m.Code, ErrResponse{Message: mType.Error()})
			if err != nil {
				return
			}
		}
	} else {
		log.Err(err).Msg("http error")
		err := c.JSON(http.StatusInternalServerError, ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		})
		if err != nil {
			return
		}
	}
}
