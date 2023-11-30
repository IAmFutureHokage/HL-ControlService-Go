package service

// import (
// 	"context"
// 	"time"

// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
// 	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
// )

// func (*HydrologyStatsService) GetDate(ctx context.Context, req *pb.GetDateRequest) (*pb.GetDateResponse, error) {
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
// 		errChan <- err
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

// 	norm := 0
// 	floodplan := 0
// 	adverse := 0
// 	dangerous := 0

// 	for _, nfad := range nfads {
// 		switch nfad.Type {
// 		case 1:
// 			norm = int(nfad.Value)
// 		case 2:
// 			floodplan = int(nfad.Value)
// 		case 3:
// 			adverse = int(nfad.Value)
// 		case 4:
// 			dangerous = int(nfad.Value)
// 		default:
// 			continue
// 		}
// 	}

// 	tx.Commit()
// 	return &pb.GetDateResponse{
// 		Data: &pb.AllNFAD{
// 			Date:       req.Date,
// 			Norm:       uint32(norm),
// 			Floodplain: uint32(floodplan),
// 			Adverse:    uint32(adverse),
// 			Dangerous:  uint32(dangerous),
// 		},
// 	}, nil
// }
