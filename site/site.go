package site

import (
	"encoding/json"
	"fmt"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages/recipe"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Site struct {
	recipreRegistry *registry.Recipe
}

func NewSite(rr *registry.Recipe) *Site {
	return &Site{
		recipreRegistry: rr,
	}
}

func (s *Site) Start(addr string) error {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			mwLog := logrus.WithFields(logrus.Fields{
				// "URI":    values.URI,
				// "status": values.Status,
				// "agent":     values.UserAgent,
				// "method":    c.Request().Method,
				"latency":   values.Latency,
				"remote_ip": values.RemoteIP,
			})

			if values.Error != nil {
				mwLog.
					WithError(values.Error).
					Error("request error")
				return values.Error
			}

			mwLog.Infof("%s %d %s", c.Request().Method, values.Status, values.URI)

			return nil
		},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())

	e.Static("/assets", "./site/public/")

	// 404
	e.GET("/*", common.UseTemplContext(pages.HandleNotFound))

	pages.InitIndexHandler(e.Group(""))
	recipe.InitEditRecipeHandler(e.Group("/recipe"), s.recipreRegistry)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return err
	}

	fmt.Println("ROUTES:")
	fmt.Println(string(data))
	fmt.Println("-------")

	// Start server
	e.Logger.Fatal(e.Start(addr))
	return nil
}

func initLogger() {
	// For a json logger
	// logrus.SetFormatter(&logrus.JSONFormatter{
	// 	FieldMap: logrus.FieldMap{
	// 		logrus.FieldKeyLevel: "severity",
	// 		logrus.FieldKeyTime:  "log_time",
	// 	},
	// })
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
}
