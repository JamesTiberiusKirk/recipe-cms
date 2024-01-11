package common

import (
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func UseBodyLogger() echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		s, _ := url.QueryUnescape(string(reqBody))
		logrus.Info("REQ BODY ", s)
		// logrus.Info("RES BODY ", string(resBody))
	})
}
