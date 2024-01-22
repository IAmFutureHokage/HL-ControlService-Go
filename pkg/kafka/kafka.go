package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaConfig struct {
	BrokerList []string `mapstructure:"broker_list"`
	Topic      string   `mapstructure:"topic"`
}

type MessageProducer interface {
	Serialize() ([]byte, error)
}

func NewKafkaProducer(config KafkaConfig) (sarama.SyncProducer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	producerConfig.Producer.Compression = sarama.CompressionSnappy
	producerConfig.Producer.Flush.Frequency = 500 * time.Millisecond
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.BrokerList, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	return producer, nil
}

func SendMessageToKafka(producer sarama.SyncProducer, topic string, messageProducer MessageProducer) error {
	messageBytes, err := messageProducer.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize message: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}
