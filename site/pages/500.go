package pages

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	conf, ok := c.Get("cfg").(config.Config)
	if !ok {
		logrus.Error("unable to get config in the custom error handler")
		return
	}

	cc := common.NewCustomContext(c)

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	logrus.Error(err)
	props := page500Props{c: cc, message: "Internal Server Error"}
	if conf.Debug {
		props.message = err.Error()
	}

	cc.TEMPL(code, page500(props))
}
