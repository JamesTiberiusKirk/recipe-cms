package common

import (
	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
)

func UseBodyLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		requestData := make(map[string]interface{})
		_ = json.NewDecoder(ctx.Request().Body).Decode(&requestData)
		log.Print(requestData)
		log.Print(ctx.Response())
		return next(ctx)
	}
}
