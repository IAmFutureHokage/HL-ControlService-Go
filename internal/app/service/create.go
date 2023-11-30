package service

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
// 	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
// 	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
// 	"github.com/google/uuid"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// func (s *HydrologyStatsService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
// 	repo := new(repository.HydrologyStatsRepository)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}

// 	prevErrChan := make(chan error, 1)
// 	prevNFADChan := make(chan *model.NFAD, 1)

// 	go func() {
// 		defer close(prevErrChan)
// 		defer close(prevNFADChan)

// 		prevNFAD, err := repo.GetActiveByPostCodeAndType(tx, int(req.GetPostCode()), byte(req.GetType()))
// 		prevErrChan <- err
// 		if err == nil {
// 			prevNFADChan <- prevNFAD
// 		}
// 	}()

// 	var prevNFAD *model.NFAD

// 	select {
// 	case err = <-prevErrChan:
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 		prevNFAD = <-prevNFADChan
// 	case <-ctx.Done():
// 		tx.Rollback()
// 		return nil, ctx.Err()
// 	}

// 	newNFAD := model.NFAD{
// 		ID:        uuid.New().String(),
// 		PostCode:  req.GetPostCode(),
// 		Type:      model.ControlType(req.GetType()),
// 		DateStart: req.GetDateStart().AsTime().Truncate(24 * time.Hour),
// 		Value:     req.GetValue(),
// 	}

// 	if prevNFAD != nil {
// 		if !newNFAD.DateStart.After(prevNFAD.DateStart) {
// 			tx.Rollback()
// 			return nil, fmt.Errorf("new NFAD's start date must be at least one day after the previous NFAD's start date")
// 		}
// 		newNFAD.PrevID = prevNFAD.ID
// 		prevNFAD.NextID = newNFAD.ID

// 		updateErrChan := make(chan error, 1)

// 		go func() {
// 			defer close(updateErrChan)
// 			updateErrChan <- repo.Update(tx, *prevNFAD)
// 		}()

// 		select {
// 		case err = <-updateErrChan:
// 			if err != nil {
// 				tx.Rollback()
// 				return nil, err
// 			}
// 		case <-ctx.Done():
// 			tx.Rollback()
// 			return nil, ctx.Err()
// 		}
// 	}

// 	createErrChan := make(chan error, 1)

// 	go func() {
// 		defer close(createErrChan)
// 		createErrChan <- repo.Create(tx, newNFAD)
// 	}()

// 	select {
// 	case err = <-createErrChan:
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	case <-ctx.Done():
// 		tx.Rollback()
// 		return nil, ctx.Err()
// 	}

// 	tx.Commit()

// 	return &pb.CreateResponse{
// 		Nfad: &pb.NFAD{
// 			Id:        newNFAD.ID,
// 			PostCode:  newNFAD.PostCode,
// 			Type:      pb.ControlType(newNFAD.Type),
// 			DateStart: timestamppb.New(newNFAD.DateStart),
// 			PrevId:    newNFAD.PrevID,
// 			Value:     newNFAD.Value,
// 		},
// 	}, nil
// }
