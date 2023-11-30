package service

// import (
// 	"context"
// 	"sort"
// 	"time"

// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
// 	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
// )

// func (*HydrologyStatsService) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
// 	repo := new(repository.HydrologyStatsRepository)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}

// 	errChan := make(chan error, 1)
// 	dataChan := make(chan []*model.NFAD, 1)

// 	go func() {
// 		defer close(errChan)
// 		defer close(dataChan)

// 		nfads, err := repo.GetByPostCodeAndDate(tx, int(req.PostCode), req.Date.AsTime().Truncate(24*time.Hour))
// 		if err == nil {
// 			dataChan <- nfads
// 		}
// 	}()

// 	var nfads []*model.NFAD

// 	select {
// 	case err = <-errChan:
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 		nfads = <-dataChan
// 	case <-ctx.Done():
// 		tx.Rollback()
// 		return nil, ctx.Err()
// 	}

// 	sort.Slice(nfads, func(i, j int) bool {
// 		return nfads[i].Value < nfads[j].Value
// 	})

// 	desiredType := 0
// 	for i := len(nfads) - 1; i >= 0; i-- {
// 		if nfads[i].Value < req.Value {
// 			desiredType = int(nfads[i].Type)
// 			break
// 		}
// 	}

// 	tx.Commit()

// 	return &pb.CheckValueResponse{
// 		Excess: uint32(desiredType),
// 	}, nil
// }
