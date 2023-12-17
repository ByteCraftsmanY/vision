package main

import (
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
	"vision/config"
	"vision/db"
	"vision/logger"
	"vision/server"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			config.Init,
			logger.Init,
			db.Init,
			server.Init,
		),
		fx.Invoke(func(db *gocqlx.Session) {}),
		fx.Invoke(func(app *http.Server) {}),
	).Run()
}
