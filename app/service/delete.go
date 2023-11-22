package service

import (
	"context"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/proto"
)

func (*ServerContext) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	repo := new(repository.RepositoryContext)

	tx, err := repo.BeginTransaction()
	if err != nil {
		return nil, err
	}

	currentNfdaErrChan := make(chan error, 1)
	currentNfadChan := make(chan *model.NFAD, 1)

	go func() {
		defer close(currentNfdaErrChan)
		defer close(currentNfadChan)

		nfad, err := repo.GetById(tx, req.Id)
		currentNfdaErrChan <- err
		if err == nil {
			currentNfadChan <- nfad
		}

	}()

	var currentNFAD *model.NFAD

	select {
	case err = <-currentNfdaErrChan:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		currentNFAD = <-currentNfadChan
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	}

	var prevNFAD, nextNFAD *model.NFAD

	if currentNFAD.PrevID != "" {

		prevErrChan := make(chan error, 1)
		prevChan := make(chan *model.NFAD, 1)

		go func() {
			defer close(prevErrChan)
			defer close(prevChan)

			prevNfad, err := repo.GetById(tx, currentNFAD.PrevID)
			prevErrChan <- err
			if err == nil {
				prevChan <- prevNfad
			}
		}()

		select {
		case err = <-prevErrChan:
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			prevNFAD = <-prevChan
		case <-ctx.Done():
			tx.Rollback()
			return nil, ctx.Err()
		}
	}

	if currentNFAD.NextID != "" {

		nextErrChan := make(chan error, 1)
		nextChan := make(chan *model.NFAD, 1)

		go func() {
			defer close(nextErrChan)
			defer close(nextChan)

			nextNfad, err := repo.GetById(tx, currentNFAD.NextID)
			nextErrChan <- err
			if err == nil {
				nextChan <- nextNfad
			}
		}()

		select {
		case err = <-nextErrChan:
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			nextNFAD = <-nextChan
		case <-ctx.Done():
			tx.Rollback()
			return nil, ctx.Err()
		}
	}

	if nextNFAD == nil && prevNFAD != nil {
		prevNFAD.NextID = ""
	}

	if nextNFAD != nil && prevNFAD == nil {
		nextNFAD.PrevID = ""
		nextNFAD.DateStart = currentNFAD.DateStart
	}

	if nextNFAD != nil && prevNFAD != nil {
		prevNFAD.NextID = nextNFAD.ID
		nextNFAD.PrevID = prevNFAD.ID
		nextNFAD.DateStart = currentNFAD.DateStart
	}

	if prevNFAD != nil {

		prevUpdateErrChan := make(chan error, 1)

		go func() {
			defer close(prevUpdateErrChan)

			err := repo.Update(tx, *prevNFAD)
			prevUpdateErrChan <- err
		}()

		select {
		case err = <-prevUpdateErrChan:
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		case <-ctx.Done():
			tx.Rollback()
			return nil, ctx.Err()
		}

	}

	if nextNFAD != nil {

		nextUpdateErrChan := make(chan error, 1)

		go func() {
			defer close(nextUpdateErrChan)

			err := repo.Update(tx, *nextNFAD)
			nextUpdateErrChan <- err
		}()

		select {
		case err = <-nextUpdateErrChan:
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		case <-ctx.Done():
			tx.Rollback()
			return nil, ctx.Err()
		}

	}

	deleteErrChan := make(chan error, 1)

	go func() {
		defer close(deleteErrChan)

		err := repo.Delete(tx, currentNFAD.ID)
		deleteErrChan <- err
	}()

	select {
	case err = <-deleteErrChan:
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	case <-ctx.Done():
		tx.Rollback()
		return nil, ctx.Err()
	}

	tx.Commit()

	return &pb.DeleteResponse{
		Success: true,
	}, nil
}
