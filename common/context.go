package common

import (
	"net/url"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
}

func NewCustomContext(c echo.Context) *Context {
	return &Context{c}
}

func (c *Context) TEMPL(status int, cmp templ.Component) error {
	c.Response().Status = status
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (c *Context) AddQueryAndReturnURI(key, value string) string {
	urlString := c.Request().URL.String()

	if strings.Contains(urlString, "?") {
		urlString += "&" + key + "=" + url.QueryEscape(value)
	} else {
		urlString += "?" + key + "=" + url.QueryEscape(value)
	}

	return urlString
}

type TemplHandlerFunc func(c *Context) error

func UseCustomContext(next TemplHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := NewCustomContext(c)
		return next(cc)
	}
}
