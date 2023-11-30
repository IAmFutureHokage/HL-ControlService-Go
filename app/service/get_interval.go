package service

// import (
// 	"context"
// 	"time"

// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
// 	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// func (*HydrologyStatsService) GetInterval(ctx context.Context, req *pb.GetIntervalRequest) (*pb.GetIntervalResponse, error) {
// 	repo := new(repository.HydrologyStatsRepository)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}

// 	errChan := make(chan error, 1)
// 	dataChan := make(chan []*model.NFAD, 1)

// 	startDate := req.StartDate.AsTime()
// 	endDate := req.EndDate.AsTime()
// 	numDays := int(endDate.Sub(startDate).Hours()/24) + 1

// 	go func() {
// 		defer close(errChan)
// 		defer close(dataChan)

// 		nfads, err := repo.GetByDateRange(tx, int(req.PostCode), startDate, endDate)
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

// 	allNfads := make([]*pb.AllNFAD, numDays)

// 	for i := 0; i < numDays; i++ {
// 		currentDay := startDate.AddDate(0, 0, i)
// 		currentDayEnd := currentDay.AddDate(0, 0, 1)

// 		allNFAD := &pb.AllNFAD{
// 			Date:       timestamppb.New(currentDay),
// 			Norm:       0,
// 			Floodplain: 0,
// 			Adverse:    0,
// 			Dangerous:  0,
// 		}

// 		for _, nfad := range nfads {
// 			if nfad.DateStart.Before(currentDayEnd) && (nfad.NextID == "" || isNextDateAfter(nfad.NextID, nfads, currentDay)) {
// 				switch nfad.Type {
// 				case 1:
// 					allNFAD.Norm = nfad.Value
// 				case 2:
// 					allNFAD.Floodplain = nfad.Value
// 				case 3:
// 					allNFAD.Adverse = nfad.Value
// 				case 4:
// 					allNFAD.Dangerous = nfad.Value
// 				}
// 			}
// 		}

// 		allNfads[i] = allNFAD
// 	}

// 	tx.Commit()
// 	return &pb.GetIntervalResponse{
// 		Data: allNfads,
// 	}, nil
// }

// func isNextDateAfter(nextID string, nfads []*model.NFAD, date time.Time) bool {
// 	for _, nfad := range nfads {
// 		if nfad.ID == nextID {
// 			return nfad.DateStart.After(date)
// 		}
// 	}
// 	return false
// }
