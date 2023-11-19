package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
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

	err = <-status
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	prevNFAD := <-prevNFADcn

	newNFAD := model.NFAD{
		ID:        uuid.New().String(),
		PostCode:  req.GetPostCode(),
		Type:      model.ControlType(req.GetType()),
		DateStart: req.GetDateStart().AsTime().Truncate(24 * time.Hour),
		Value:     req.GetValue(),
	}

	if prevNFAD != nil {
		if !newNFAD.DateStart.After(prevNFAD.DateStart) {
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

func (*ServerContext) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	status := make(chan error)
	data := make(chan *model.NFAD)

	go repo.GetById(tx, req.Id, status, data)

	err = <-status
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	getNFDA := <-data

	var prevNFAD, nextNFAD *model.NFAD

	if getNFDA.PrevID != "" {

		status = make(chan error)
		data = make(chan *model.NFAD)

		go repo.GetById(tx, getNFDA.PrevID, status, data)

		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		prevNFAD = <-data
	}

	if getNFDA.NextID != "" {

		status = make(chan error)
		data = make(chan *model.NFAD)

		go repo.GetById(tx, getNFDA.NextID, status, data)

		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		nextNFAD = <-data
	}

	if nextNFAD == nil && prevNFAD != nil {
		prevNFAD.NextID = ""
	}

	if nextNFAD != nil && prevNFAD == nil {
		nextNFAD.PrevID = ""
		nextNFAD.DateStart = getNFDA.DateStart
	}

	if nextNFAD != nil && prevNFAD != nil {
		prevNFAD.NextID = nextNFAD.ID
		nextNFAD.PrevID = prevNFAD.ID
		nextNFAD.DateStart = getNFDA.DateStart
	}

	if prevNFAD != nil {
		status = make(chan error)
		go repo.Update(tx, *prevNFAD, status)

		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if nextNFAD != nil {
		status = make(chan error)
		go repo.Update(tx, *nextNFAD, status)

		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	status = make(chan error)

	go repo.Delete(tx, getNFDA.ID, status)

	err = <-status
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}

func (*ServerContext) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}
	if len(req.Data) < 1 {
		return nil, errors.New("no data provided")
	}

	if len(req.Data) != 1 {

		for i := 0; i < len(req.Data)-1; i++ {
			current := req.Data[i]
			next := req.Data[i+1]

			if current.PostCode != next.PostCode ||
				current.Type != next.Type ||
				current.Id != next.PrevId ||
				!next.DateStart.AsTime().Truncate(24*time.Hour).After(current.DateStart.AsTime().Truncate(24*time.Hour)) {
				return nil, errors.New("bad data")
			}
		}
	}

	status := make(chan error)
	close(status)

	for _, pbNFAD := range req.Data {

		nfad := model.NFAD{
			ID:        pbNFAD.Id,
			PostCode:  pbNFAD.PostCode,
			Type:      model.ControlType(pbNFAD.Type),
			DateStart: pbNFAD.DateStart.AsTime().Truncate(24 * time.Hour),
			Value:     pbNFAD.Value,
			PrevID:    pbNFAD.PrevId,
			NextID:    pbNFAD.NextId,
		}

		status = make(chan error)
		go repo.Update(tx, nfad, status)
		err = <-status
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	return &pb.UpdateResponse{
		Data: req.Data,
	}, nil
}

func (*ServerContext) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	const pageSize = 50

	dataChan := make(chan []*model.NFAD)
	statusChan := make(chan error)
	totalPagesChan := make(chan int)

	go repo.GetByPostCodeAndType(tx, int(req.PostCode), byte(req.Type), int(req.Page), pageSize, statusChan, dataChan, totalPagesChan)
	err = <-statusChan
	if err != nil {
		return nil, err
	}

	maxPages := <-totalPagesChan
	nfads := <-dataChan

	pbNfads := make([]*pb.NFAD, len(nfads))
	for i, nfad := range nfads {
		pbNfads[i] = &pb.NFAD{
			Id:        nfad.ID,
			PostCode:  nfad.PostCode,
			Type:      pb.ControlType(nfad.Type.ToByte()),
			DateStart: timestamppb.New(nfad.DateStart),
			PrevId:    nfad.PrevID,
			NextId:    nfad.NextID,
			Value:     nfad.Value,
		}
	}

	tx.Commit()

	return &pb.GetResponse{
		Page:    req.GetPage(),
		MaxPage: uint32(maxPages),
		Data:    pbNfads,
	}, nil
}

func (*ServerContext) CheckValue(ctx context.Context, req *pb.CheckValueRequest) (*pb.CheckValueResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []*model.NFAD)
	statusChan := make(chan error)

	go repo.GetByPostCodeAndDate(tx, int(req.PostCode), req.Date.AsTime().Truncate(24*time.Hour), statusChan, dataChan)

	err = <-statusChan
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	nfads := <-dataChan

	sort.Slice(nfads, func(i, j int) bool {
		return nfads[i].Value < nfads[j].Value
	})

	desiredType := 0
	for i := len(nfads) - 1; i >= 0; i-- {
		if nfads[i].Value < req.Value {
			desiredType = int(nfads[i].Type)
			break
		}
	}

	tx.Commit()

	return &pb.CheckValueResponse{
		Excess: uint32(desiredType),
	}, nil
}

func (*ServerContext) GetDate(ctx context.Context, req *pb.GetDateRequest) (*pb.GetDateResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []*model.NFAD)
	statusChan := make(chan error)

	go repo.GetByPostCodeAndDate(tx, int(req.PostCode), req.Date.AsTime().Truncate(24*time.Hour), statusChan, dataChan)

	err = <-statusChan
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	nfads := <-dataChan

	norm := 0
	floodplan := 0
	adverse := 0
	dangerous := 0

	for _, nfad := range nfads {
		switch nfad.Type {
		case 1:
			norm = int(nfad.Value)
		case 2:
			floodplan = int(nfad.Value)
		case 3:
			adverse = int(nfad.Value)
		case 4:
			dangerous = int(nfad.Value)
		default:
			continue
		}
	}

	tx.Commit()
	return &pb.GetDateResponse{
		Data: &pb.AllNFAD{
			Date:       req.Date,
			Norm:       uint32(norm),
			Floodplain: uint32(floodplan),
			Adverse:    uint32(adverse),
			Dangerous:  uint32(dangerous),
		},
	}, nil
}

func (*ServerContext) GetInterval(ctx context.Context, req *pb.GetIntervalRequest) (*pb.GetIntervalResponse, error) {
	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []*model.NFAD)
	statusChan := make(chan error)

	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()
	numDays := int(endDate.Sub(startDate).Hours()/24) + 1

	go repo.GetByDateRange(tx, int(req.PostCode), startDate, endDate, statusChan, dataChan)

	err = <-statusChan
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	nfads := <-dataChan

	allNfads := make([]*pb.AllNFAD, numDays)

	for i := 0; i < numDays; i++ {
		currentDay := startDate.AddDate(0, 0, i)
		currentDayEnd := currentDay.AddDate(0, 0, 1)

		allNFAD := &pb.AllNFAD{
			Date:       timestamppb.New(currentDay),
			Norm:       0,
			Floodplain: 0,
			Adverse:    0,
			Dangerous:  0,
		}

		for _, nfad := range nfads {
			if nfad.DateStart.Before(currentDayEnd) && (nfad.NextID == "" || isNextDateAfter(nfad.NextID, nfads, currentDay)) {
				switch nfad.Type {
				case 1:
					allNFAD.Norm = nfad.Value
				case 2:
					allNFAD.Floodplain = nfad.Value
				case 3:
					allNFAD.Adverse = nfad.Value
				case 4:
					allNFAD.Dangerous = nfad.Value
				}
			}
		}

		allNfads[i] = allNFAD
	}

	tx.Commit()
	return &pb.GetIntervalResponse{
		Data: allNfads,
	}, nil
}

func isNextDateAfter(nextID string, nfads []*model.NFAD, date time.Time) bool {
	for _, nfad := range nfads {
		if nfad.ID == nextID {
			return nfad.DateStart.After(date)
		}
	}
	return false
}
