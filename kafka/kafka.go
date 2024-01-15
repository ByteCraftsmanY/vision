package kafka

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"
	"vision/config"
)

var client *kgo.Client

func Init(lc fx.Lifecycle, config *config.Configurations) *kgo.Client {
	err := error(nil)
	kafkaConfig := config.Kafka

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			client, err = kgo.NewClient(
				kgo.DefaultProduceTopic(kafkaConfig.DefaultTopic),
				kgo.SeedBrokers(kafkaConfig.Brokers...),
				kgo.RequiredAcks(kgo.AllISRAcks()),
				//kgo.WithLogger(zap.L()),
			)
			return err
		},
		OnStop: func(ctx context.Context) error {
			client.Close()
			return err
		},
	})
	return client
}

func GetClient() *kgo.Client {
	return client
}
