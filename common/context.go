package common

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
}

func (c *Context) TEMPL(status int, cmp templ.Component) error {
	c.Response().Status = status
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

type TemplHandlerFunc func(c *Context) error

func UseCustomContext(next TemplHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{c}
		return next(cc)
	}
}
