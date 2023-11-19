package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/service"
	"github.com/IAmFutureHokage/HL-ControlService-Go/database"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func main() {
	fmt.Println("gRPC server running ...")

	port := viper.GetInt("server.port")
	if port == 0 {
		log.Fatal("Server port is not set in the config file")
	}

	_, _ = database.OpenDB()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHydrologyControlServiceServer(s, &service.ServerContext{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
