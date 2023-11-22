package service

import (
	// "context"
	// "errors"

	// "time"

	// "github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	// "github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

type ServerContext struct {
	pb.UnimplementedHydrologyControlServiceServer
}

// func (*ServerContext) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
// 	repo := new(repository.RepositoryContext)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(req.Data) < 1 {
// 		return nil, errors.New("no data provided")
// 	}

// 	if len(req.Data) != 1 {

// 		for i := 0; i < len(req.Data)-1; i++ {
// 			current := req.Data[i]
// 			next := req.Data[i+1]

// 			if current.PostCode != next.PostCode ||
// 				current.Type != next.Type ||
// 				current.Id != next.PrevId ||
// 				!next.DateStart.AsTime().Truncate(24*time.Hour).After(current.DateStart.AsTime().Truncate(24*time.Hour)) {
// 				return nil, errors.New("bad data")
// 			}
// 		}
// 	}

// 	status := make(chan error)
// 	close(status)

// 	for _, pbNFAD := range req.Data {

// 		nfad := model.NFAD{
// 			ID:        pbNFAD.Id,
// 			PostCode:  pbNFAD.PostCode,
// 			Type:      model.ControlType(pbNFAD.Type),
// 			DateStart: pbNFAD.DateStart.AsTime().Truncate(24 * time.Hour),
// 			Value:     pbNFAD.Value,
// 			PrevID:    pbNFAD.PrevId,
// 			NextID:    pbNFAD.NextId,
// 		}

// 		status = make(chan error)
// 		go repo.Update(tx, nfad, status)
// 		err = <-status
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	}

// 	tx.Commit()

// 	return &pb.UpdateResponse{
// 		Data: req.Data,
// 	}, nil
// }
