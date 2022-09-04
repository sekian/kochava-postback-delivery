package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/segmentio/kafka-go"
)

const (
	topic      = "postback"
	broker_uri = "localhost:9092"
)

func consume() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker_uri},
		Topic:   topic,
		GroupID: topic,
	})
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		fmt.Println(reflect.TypeOf(msg.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func main() {
	consume()
}
