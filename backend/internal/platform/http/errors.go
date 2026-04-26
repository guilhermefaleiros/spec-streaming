package http

import "github.com/labstack/echo/v4"

func NewErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		c.JSON(500, map[string]string{"error": err.Error()})
	}
}
