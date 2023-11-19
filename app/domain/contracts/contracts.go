package contracts

import (
	"context"
	"time"

	model "github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error)
	Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error)
	Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error)
	Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error)
	CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error)
	GetDate(ctx context.Context, req *pb.GetDateRequest) (*pb.GetDateResponse, error)
	GetInterval(ctx context.Context, req *pb.GetIntervalRequest) (*pb.GetIntervalResponse, error)
}

type Repository interface {
	Create(tx *gorm.DB, data model.NFAD, status chan error)
	Delete(tx *gorm.DB, id string, status chan error)
	Update(tx *gorm.DB, data model.NFAD, status chan error)
	GetById(tx *gorm.DB, id string, status chan error, data chan *model.NFAD)
	GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int, status chan error, data chan []*model.NFAD, totalPages chan int)
	GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, status chan error, data chan *model.NFAD)
	GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time, status chan error, data chan []*model.NFAD)
	GetByDateRange(tx *gorm.DB, postCode int, startDate time.Time, endDate time.Time, status chan error, data chan []*model.NFAD)
}
