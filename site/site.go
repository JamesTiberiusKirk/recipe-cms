package site

import (
	"log"
	"os"
	"strings"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages/auth"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages/playground"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages/recipe"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Site struct {
	userRegistry    registry.IUser
	recipreRegistry registry.IRecipe
	config          config.Config
	sessions        *session.Manager
}

func NewSite(conf config.Config, rr registry.IRecipe, ur registry.IUser) *Site {
	return &Site{
		config:          conf,
		recipreRegistry: rr,
		userRegistry:    ur,
		sessions:        session.New(),
	}
}

func (s *Site) Start(addr string) error {
	_, err := os.Stat(s.config.Volume)
	if os.IsNotExist(err) {
		log.Fatalf("mounted volume does not exist, path: %s, err: %s", s.config.Volume, err.Error())
	}
	if err != nil {
		log.Fatalf("error getting mounted volume stat: %s, err: %s", s.config.Volume, err.Error())
	}

	// Echo instance
	e := echo.New()

	e.Use(
		inject(s.config, s.sessions),
	)

	// Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {

			// logging excludes
			if strings.Contains(values.URI, "/images") {
				return nil
			}

			if strings.Contains(values.URI, "/assets") {
				return nil
			}

			if strings.Contains(values.URI, "/favicon.ico") {
				return nil
			}

			if strings.Contains(values.URI, "/_templ/reload/events") {
				return nil
			}

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

	e.HTTPErrorHandler = pages.CustomHTTPErrorHandler

	e.Static("/images", s.config.Volume)
	e.Static("/assets", "./site/public/")

	// 404
	e.RouteNotFound("/*", common.UseCustomContext(pages.HandleNotFound))

	if s.config.Debug {
		e.GET("/ready", pages.Ready)
	}

	auth.InitAuthHandler(e.Group("/auth"), s.sessions, s.userRegistry)
	pages.InitIndexHandler(e.Group(""))
	recipe.InitRecipeHandler(e.Group("/recipe"), s.recipreRegistry, s.sessions)
	recipe.InitRecipesHandler(e.Group("/recipes"), s.recipreRegistry)
	pages.InitMarkdownRenderer(e.Group(""))
	playground.InitTestRoute(e.Group("/pg"))

	// data, err := json.MarshalIndent(e.Routes(), "", "  ")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("ROUTES:")
	// fmt.Println(string(data))
	// fmt.Println("-------")

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

func inject(cfg config.Config, m *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("cfg", cfg)
			c.Set("session", m)
			return next(c)
		}
	}
}
