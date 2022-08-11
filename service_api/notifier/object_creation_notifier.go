package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

var ObjectCreationNotificationChannel = make(chan interface{}, 50)

type NotifierFunc func(chan interface{})

func getKafkaWriter() *kafka.Writer {
	viper.BindEnv("KAFKA_URL")
	viper.BindEnv("OBJECT_CREATION_TOPIC_NAME")
	kafka_url := viper.GetString("KAFKA_URL")
	topic := viper.GetString("OBJECT_CREATION_TOPIC_NAME")
	fmt.Printf("Creating kafka writer with url '%s' and topic '%s'", kafka_url, topic)
	return &kafka.Writer{
		Addr:     kafka.TCP(kafka_url),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

var kafkaWriter = getKafkaWriter()

func KafkaNotifier(c chan interface{}) {
	for val := range c {
		msg := make(map[string]interface{})
		msg["type"] = reflect.TypeOf(val).Name()
		msg["value"] = val
		r, _ := json.Marshal(msg)
		fmt.Println(string(r))

		kafka_msg := kafka.Message{
			Key:   []byte(reflect.TypeOf(val).Name()),
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
