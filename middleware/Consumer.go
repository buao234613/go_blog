package middleware

//import (
//	"WhiteBlog/common"
//	"WhiteBlog/models"
//	"context"
//	"encoding/json"
//	"github.com/olivere/elastic/v7"
//	"github.com/segmentio/kafka-go"
//	"log"
//	"runtime/debug"
//	"strconv"
//	"sync"
//)
//
//type Consumer struct {
//	reader *kafka.Reader
//}
//
//var (
//	onceConsumer sync.Once
//	consumer     *Consumer
//)
//
//func GetConsumer() *Consumer {
//	onceConsumer.Do(initConsumer)
//	return consumer
//}
//func initConsumer() {
//	// Kafka 消费者配置
//	config := kafka.ReaderConfig{
//		Brokers:  []string{"localhost:9092"}, // 替换为 Kafka 的实际地址
//		Topic:    "article_events",
//		GroupID:  "articles_group",
//		MinBytes: 10e3, // 10KB
//		MaxBytes: 10e6, // 10MB
//	}
//	// 初始化消费者
//	reader := kafka.NewReader(config)
//	consumer = &Consumer{
//		reader: reader,
//	}
//}
//func (c *Consumer) ConsumeEvent(client *elastic.Client, event Event) {
//	for {
//		// 读取消息
//		msg, err := c.reader.ReadMessage(context.Background())
//		if err != nil {
//			log.Printf("Error reading message from Kafka: %v", err)
//			debug.PrintStack()
//			continue
//		}
//
//		err = json.Unmarshal(msg.Value, &event)
//		if err != nil {
//			log.Printf("Error unmarshaling message: %v", err)
//			continue
//		}
//		err = handleEvent(client, event)
//		if err != nil {
//			log.Println("handle event error:", err)
//		}
//	}
//}
//func handleEvent(client *elastic.Client, event Event) error {
//	switch event.Topic {
//	case common.Create:
//		return Create(client, event.Article)
//	case common.Update:
//		return Update(client, event.Article)
//	case common.Delete:
//		return Delete(client, event.Article)
//	}
//	return nil
//}
//
//
