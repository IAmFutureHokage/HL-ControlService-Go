package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/migrations"
	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/repository"
	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/service"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/internal/proto"
	"github.com/IAmFutureHokage/HL-ControlService-Go/pkg/database"
	"github.com/IAmFutureHokage/HL-ControlService-Go/pkg/kafka"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var dbConfig database.Config
var kafkaConfig kafka.KafkaConfig
var kafkaProducer sarama.SyncProducer

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(env)
	viper.AddConfigPath("../../config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	dbConfig = database.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		PoolSize: viper.GetInt("database.poolsize"),
	}

	kafkaConfig = kafka.KafkaConfig{
		BrokerList: viper.GetStringSlice("kafka.broker_list"),
		Topic:      viper.GetString("kafka.topic"),
	}

	var err error
	kafkaProducer, err = kafka.NewKafkaProducer(kafkaConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
}

func main() {

	fmt.Println("gRPC server running ...")

	port := viper.GetInt("server.port")
	if port == 0 {
		log.Fatal("Server port is not set in the config file")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	dbPool, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer database.CloseDB(dbPool)

	if _, err := dbPool.Exec(context.Background(), migrations.CreateTableControlValue); err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	repo := repository.NewHydrologyStatsRepository(dbPool)
	hydrologyStatsService := service.NewHydrologyStatsService(repo)
	kafkaMessageService := service.NewKafkaMessageService(repo)

	go func() {
		kafka.SubscribeToTopic(kafkaConfig, kafkaMessageService)
	}()

	s := grpc.NewServer()
	pb.RegisterHydrologyStatsServiceServer(s, hydrologyStatsService)

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
