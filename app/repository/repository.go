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
		close(status)
		return
	}

	res := db.Create(&data)

	if res.RowsAffected == 0 {
		err = errors.New("failed create")
		status <- err
		close(status)
		return
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) Delete(id string, status chan error) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		return
	}

	var nfad model.NFAD
	res := db.Where("id=?", id).Delete(&nfad)
	if res.RowsAffected == 0 {
		status <- errors.New("not found")
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) Update(data model.NFAD, status chan error) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		return
	}

	res := db.Model(&model.NFAD{}).Where("id = ?", data.ID).Updates(data)

	if res.Error != nil {
		status <- res.Error
		close(status)
		return
	}

	if res.RowsAffected == 0 {
		status <- errors.New("no rows affected")
		close(status)
		return
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) GetById(id string, status chan error, data chan *model.NFAD) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		close(data)
		return
	}

	var nfad model.NFAD
	res := db.First(&nfad, "id = ?", id)

	if res.Error != nil {
		status <- res.Error
		close(status)
		close(data)
		return
	}

	if res.RowsAffected == 0 {
		status <- errors.New("no rows affected")
		close(status)
		close(data)
		return
	}

	data <- &nfad
	close(data)
	status <- nil
	close(status)
}
