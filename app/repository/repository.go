package repository

import (
	"errors"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/app/domain/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/database"
	"gorm.io/gorm"
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

	updateData := map[string]interface{}{
		"PostCode":  data.PostCode,
		"Type":      data.Type,
		"DateStart": data.DateStart,
		"PrevID":    data.PrevID,
		"NextID":    data.NextID,
		"Value":     data.Value,
	}

	res := db.Model(&model.NFAD{}).Where("id = ?", data.ID).Updates(updateData)

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

func (r RepositoryContext) GetAllByPostCodeAndType(postCode int, typeNfad byte, status chan error, data chan []*model.NFAD) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		close(data)
		return
	}

	var nfads []*model.NFAD
	res := db.Where("post_code = ? AND type = ?", postCode, typeNfad).Find(&nfads)

	if res.Error != nil {
		status <- res.Error
		close(status)
		close(data)
		return
	}

	if len(nfads) == 0 {
		status <- errors.New("slice is empty")
		close(status)
		close(data)
		return
	}

	data <- nfads
	close(data)
	status <- nil
	close(status)
}

func (r RepositoryContext) GetActiveByPostCodeAndType(postCode int, typeNfad byte, status chan error, data chan *model.NFAD) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		close(data)
		return
	}

	var nfad model.NFAD
	res := db.Where("post_code = ? AND type = ? AND next_id = ?", postCode, typeNfad, "").First(&nfad)

	if res.RowsAffected == 0 {
		data <- nil
		close(data)
		status <- nil
		close(status)
		return
	}

	if res.Error != nil {
		status <- res.Error
		close(status)
		close(data)
		return
	}

	data <- &nfad
	close(data)
	status <- nil
	close(status)
}

func (r RepositoryContext) GetByPostCodeAndDate(postCode int, date time.Time, status chan error, data chan []*model.NFAD) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		close(data)
		return
	}

	date = date.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	for _, typeNfad := range []byte{1, 2, 3, 4} {
		var nfad model.NFAD
		res := db.Where("post_code = ? AND type = ? AND date_start <= ?", postCode, typeNfad, date).Order("date_start desc").First(&nfad)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				continue
			}
			status <- res.Error
			close(status)
			close(data)
			return
		}

		nfads = append(nfads, &nfad)
	}

	if len(nfads) == 0 {
		status <- errors.New("no records found for the specified post code and date")
		close(status)
		close(data)
		return
	}

	data <- nfads
	close(data)
	status <- nil
	close(status)
}

func (r RepositoryContext) GetByDateRange(postCode int, startDate time.Time, endDate time.Time, status chan error, data chan []*model.NFAD) {
	db, err := database.OpenDB()
	if err != nil {
		status <- err
		close(status)
		close(data)
		return
	}

	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	for _, typeNfad := range []byte{1, 2, 3, 4} {
		var nfad model.NFAD
		res := db.Where("post_code = ? AND type = ? AND date_start <= ?", postCode, typeNfad, endDate).Order("date_start DESC").First(&nfad)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				continue
			}
			status <- res.Error
			close(status)
			close(data)
			return
		}

		if nfad.DateStart.Before(startDate) {
			nfads = append(nfads, &nfad)
		}
	}

	if len(nfads) == 0 {
		status <- errors.New("no records found for the specified range")
		close(status)
		close(data)
		return
	}

	data <- nfads
	close(data)
	status <- nil
	close(status)
}
