package middleware

import (
	"WhiteBlog/models"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

type Event struct {
	Topic   string
	Article models.Article
}

var (
	// 生产者
	producer     *kafka.Writer
	onceProducer sync.Once
)

func GetProducer() *kafka.Writer {
	onceProducer.Do(initProducer)
	return producer
}

func initProducer() {
	config := kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"}, // 替换为 Kafka 的实际地址
		Topic:    "default_topic",            // 默认主题
		Balancer: &kafka.LeastBytes{},
	}
	producer = kafka.NewWriter(config)
}
func SendEvent(event Event) error {
	if producer == nil {
		initProducer()
	}
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		log.Println("send event:", event, " error :", err)
	}
	// 创建 Kafka 消息
	message := kafka.Message{
		Key:   []byte(event.Topic), // 使用主题作为消息的键
		Value: jsonEvent,
	}
	// 发送消息
	ctx := context.Background()
	err = producer.WriteMessages(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
