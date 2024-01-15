package server

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"vision/config"
)

type Params struct {
	fx.In

	Config *config.Configurations
	Logger *zap.Logger
}

func Init(lc fx.Lifecycle, params Params) *http.Server {
	r := NewRouter(&params)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.Config.Server.Port),
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					zap.L().Panic("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
