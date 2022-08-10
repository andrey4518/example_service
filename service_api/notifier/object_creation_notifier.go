package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	kafka "github.com/segmentio/kafka-go"
)

var ObjectCreationNotificationChannel = make(chan interface{}, 50)

type NotifierFunc func(chan interface{})

var kafka_url = "example_service_kafka_1:9092"
var topic = "test-topic"

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

var kafkaWriter = getKafkaWriter(kafka_url, topic)

func KafkaNotifier(c chan interface{}) {
	for val := range c {
		msg := make(map[string]interface{})
		msg["type"] = reflect.TypeOf(val).Name()
		msg["value"] = val
		r, _ := json.Marshal(msg)
		fmt.Println(string(r))

		kafka_msg := kafka.Message{
			Key:   []byte("test_key"),
			Value: r,
		}

		err := kafkaWriter.WriteMessages(context.Background(), kafka_msg)

		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}
}

func CreateObjectCreationNotifierFunc() NotifierFunc {
	return KafkaNotifier
}
