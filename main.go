package main

import (
	"github.com/scylladb/gocqlx/v2"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"net/http"
	"vision/config"
	"vision/db"
	"vision/kafka"
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
			kafka.Init,
			server.Init,
		),
		fx.Invoke(func(db *gocqlx.Session) {}),
		fx.Invoke(func(client *kgo.Client) {}),
		fx.Invoke(func(app *http.Server) {}),
	).Run()
}
