package repository

import (
	"errors"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/database"
)

type RepositoryContext struct {
}

func (r RepositoryContext) Create(data model.NFAD, status chan error) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
	}

	res := db.Create(&data)

	if res.RowsAffected == 0 {
		err = errors.New("failed create movie")
		status <- err
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) Delete(id string, status chan error) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
	}

	var nfad model.NFAD
	res := db.Where("id=?", id).Delete(&nfad)
	if res.RowsAffected == 0 {
		status <- errors.New("movies not found")
	}

	status <- nil
	close(status)
}
