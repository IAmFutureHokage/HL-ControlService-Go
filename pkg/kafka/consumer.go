package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type MessageConsumer interface {
	ProcessMessage([]byte) error
}

func SubscribeToTopic(config KafkaConfig, consumer MessageConsumer) {
	brokers := config.BrokerList
	topic := config.Topic
	groupID := "my-group"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.LastOffset,
		MaxWait:     1 * time.Second,
	})

	defer r.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			log.Println("Received signal. Closing consumer.")
			return
		default:
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				break
			}

			if err := consumer.ProcessMessage(m.Value); err != nil {
				log.Printf("Error processing Kafka message: %v\n", err)
			}
		}
	}
}
