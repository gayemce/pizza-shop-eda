package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gayemce/pizza-shop-eda/order-service/config"
	"github.com/gayemce/pizza-shop-eda/order-service/logger"
	"github.com/segmentio/kafka-go"
)

type IMessageConsumer interface {
	ConsumeMessage() (Message, error)
	GetReader() *kafka.Reader
	Close() error
}

type Message struct {
	Data         map[string]interface{}
	KafkaMessage kafka.Message
	Topic        string
}

type KafkaMessageConsumer struct {
	conn   *config.KafkaConnection
	Reader *kafka.Reader
}

func (kc *KafkaMessageConsumer) ConsumeMessage() (Message, error) {
	var data map[string]interface{}
	var event = Message{}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	msg, err := kc.Reader.ReadMessage(ctx)
	if err != nil {
		return event, fmt.Errorf("failed to read message from kafka: %v", err)
	}
	err = json.Unmarshal(msg.Value, &data)
	if err != nil {
		return event, fmt.Errorf("parse message error: %v", err)
	}
	event.Data = data
	event.KafkaMessage = msg

	return event, nil
}

func (kc *KafkaMessageConsumer) GetReader() *kafka.Reader {
	return kc.Reader
}

func (kc *KafkaMessageConsumer) Close() error {
	err := kc.Reader.Close()
	if err != nil {
		logger.Log(fmt.Sprintf("Error closing Kafka Reader: %v", err))
		return err
	}
	logger.Log("Kafka Reader closed successfully")
	return nil
}

func GetNewKafkaConsumer(topic, groupId string) *KafkaMessageConsumer {
	conn := config.GetNewKafkaConnection(topic, groupId)
	reader := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   topic,
			GroupID: groupId,
		},
	)
	return &KafkaMessageConsumer{
		conn:   conn,
		Reader: reader,
	}
}
