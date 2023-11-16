package services

import (
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
)

type ServerContext struct {
	pb.UnimplementedHydrologyControlServiceServer
}
