package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/model"
	"github.com/google/uuid"
)

type WaterLevelStorage interface {
	AddWaterlevel(ctx context.Context, value model.Waterlevel) error
}

type KafkaMessageService struct {
	storage WaterLevelStorage
}

func NewKafkaMessageService(storage WaterLevelStorage) *KafkaMessageService {
	return &KafkaMessageService{storage: storage}
}

func (s *KafkaMessageService) ProcessMessage(message []byte) error {

	ctx := context.Background()

	var kafkaMessage model.KafkaMessage
	if err := json.Unmarshal(message, &kafkaMessage); err != nil {
		return fmt.Errorf("failed to unmarshal Kafka message: %v", err)
	}

	waterlevel := &model.Waterlevel{
		Id:         uuid.New(),
		PostCode:   kafkaMessage.PostCode,
		Date:       kafkaMessage.Date,
		Waterlevel: kafkaMessage.WaterLevel,
	}

	fmt.Println(waterlevel)

	if err := s.storage.AddWaterlevel(ctx, *waterlevel); err != nil {
		fmt.Println("Не удалось добавить в базу")
	}

	return nil
}
