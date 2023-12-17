package db

import (
	"context"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/fx"
	"log"
	"vision/config"
)

var session gocqlx.Session

func Init(lc fx.Lifecycle, config *config.Configurations) *gocqlx.Session {
	err := error(nil)
	dbConfig := &config.Database

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			clusterConfig := gocql.NewCluster(dbConfig.Hosts...)
			clusterConfig.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
				NumRetries: 3,
				Min:        1,
				Max:        5,
			}
			clusterConfig.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
			clusterConfig.Consistency = gocql.LocalQuorum
			clusterConfig.Keyspace = dbConfig.Keyspace
			clusterConfig.Logger = log.Default()
			session, err = gocqlx.WrapSession(clusterConfig.CreateSession())
			return err
		},
		OnStop: func(ctx context.Context) error {
			session.Close()
			return nil
		},
	})
	return &session
}

func GetSession() *gocqlx.Session {
	return &session
}
