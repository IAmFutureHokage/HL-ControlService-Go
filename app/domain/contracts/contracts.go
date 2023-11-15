package contracts

import (
	"context"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/dto"
	model "github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
)

type ServiceGrpc interface {
	Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error)
	Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error)
	CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error)
	GetDate(ctx context.Context, req *pb.GetDateRequest) (*pb.GetDateResponse, error)
	GetInterval(ctx context.Context, req *pb.GetIntervalRequest) (*pb.GetIntervalResponse, error)
}

type Repository interface {
	Create(postCode int32, controlType model.ControlType, dateStart time.Time, value int32, response chan<- model.NFAD, errChan chan<- error)
	Delete(id string, response chan<- bool, errChan chan<- error)
	Update(data []model.NFAD, response chan<- []model.NFAD, errChan chan<- error)
	Get(postCode int32, controlType model.ControlType, page uint32, response chan<- dto.GetResponse, errChan chan<- error)
	CheckValue(date time.Time, postCode int32, value int32, response chan<- dto.CheckValueResponse, errChan chan<- error)
	GetDate(postCode int32, date time.Time, response chan<- dto.GetDateResponse, errChan chan<- error)
	GetInterval(postCode int32, startDate time.Time, endDate time.Time, response chan<- dto.GetIntervalResponse, errChan chan<- error)
}
