package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServerContext struct {
	pb.UnimplementedHydrologyControlServiceServer
}

func (*ServerContext) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	status := make(chan error)
	prevNFADcn := make(chan *model.NFAD)

	go repo.GetActiveByPostCodeAndType(tx, int(req.GetPostCode()), byte(req.GetType()), status, prevNFADcn)

	prevNFAD := <-prevNFADcn

	err = <-status
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	newNFAD := model.NFAD{
		ID:        uuid.New().String(),
		PostCode:  req.GetPostCode(),
		Type:      model.ControlType(req.GetType()),
		DateStart: req.GetDateStart().AsTime().Truncate(24 * time.Hour),
		Value:     req.GetValue(),
	}

	if prevNFAD != nil {
		if !newNFAD.DateStart.After(prevNFAD.DateStart.AddDate(0, 0, 1)) {
			tx.Rollback()
			return nil, fmt.Errorf("new NFAD's start date must be at least one day after the previous NFAD's start date")
		}
		newNFAD.PrevID = prevNFAD.ID
		prevNFAD.NextID = newNFAD.ID

		status = make(chan error)
		go repo.Update(tx, *prevNFAD, status)

		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	status = make(chan error)
	go repo.Create(tx, newNFAD, status)

	err = <-status
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &pb.CreateResponse{
		Nfad: &pb.NFAD{
			Id:        newNFAD.ID,
			PostCode:  newNFAD.PostCode,
			Type:      pb.ControlType(newNFAD.Type),
			DateStart: timestamppb.New(newNFAD.DateStart),
			PrevId:    newNFAD.PrevID,
			Value:     newNFAD.Value,
		},
	}, nil
}
