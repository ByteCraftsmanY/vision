package services

import (
	"context"
	"encoding/json"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
	"vision/kafka"
)

type KafkaService struct{}

func getRecord(ctx *context.Context, topic *string, msg interface{}) (*kgo.Record, error) {
	msgBytes, err := json.Marshal(msg)
	record := kgo.Record{
		Key:     []byte(*topic),
		Value:   msgBytes,
		Topic:   *topic,
		Context: *ctx,
	}
	return &record, err
}

func (s *KafkaService) ProduceMessage(ctx context.Context, topic string, msg interface{}) error {
	client := kafka.GetClient()
	record, err := getRecord(&ctx, &topic, &msg)
	if err != nil {
		return err
	}
	client.Produce(ctx, record, func(record *kgo.Record, err error) {
		if err != nil {
			zap.L().Error("Failed to write async msg in kafka", zap.Error(err))
		}
	})
	return err
}

func (s *KafkaService) ProduceMessageSync(ctx context.Context, topic string, messages ...interface{}) error {
	client := kafka.GetClient()
	records := make([]*kgo.Record, 0)
	for _, message := range messages {
		record, err := getRecord(&ctx, &topic, &message)
		if err != nil {
			return err
		}
		records = append(records, record)
	}
	results := client.ProduceSync(ctx, records...)
	for _, r := range results {
		if r.Err != nil {
			return r.Err
		}
	}
	return nil
}
