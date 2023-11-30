package service

// import (
// 	"context"

// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
// 	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// func (*HydrologyStatsService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
// 	repo := new(repository.HydrologyStatsRepository)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}

// 	const pageSize = 50

// 	errChan := make(chan error, 1)
// 	totalPagesChan := make(chan int, 1)
// 	dataChan := make(chan []*model.NFAD, 1)

// 	go func() {
// 		defer close(dataChan)
// 		defer close(totalPagesChan)
// 		defer close(errChan)

// 		total, nfads, err := repo.GetByPostCodeAndType(tx, int(req.PostCode), byte(req.Type), int(req.Page), pageSize)
// 		errChan <- err
// 		if err == nil {
// 			totalPagesChan <- total
// 			dataChan <- nfads
// 		}
// 	}()

// 	select {
// 	case err = <-errChan:
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	case <-ctx.Done():
// 		tx.Rollback()
// 		return nil, ctx.Err()
// 	}

// 	maxPages := <-totalPagesChan
// 	nfads := <-dataChan

// 	pbNfads := make([]*pb.NFAD, len(nfads))
// 	for i, nfad := range nfads {
// 		pbNfads[i] = &pb.NFAD{
// 			Id:        nfad.ID,
// 			PostCode:  nfad.PostCode,
// 			Type:      pb.ControlType(nfad.Type.ToByte()),
// 			DateStart: timestamppb.New(nfad.DateStart),
// 			PrevId:    nfad.PrevID,
// 			NextId:    nfad.NextID,
// 			Value:     nfad.Value,
// 		}
// 	}

// 	tx.Commit()

// 	return &pb.GetResponse{
// 		Page:    req.GetPage(),
// 		MaxPage: uint32(maxPages),
// 		Data:    pbNfads,
// 	}, nil
// }
