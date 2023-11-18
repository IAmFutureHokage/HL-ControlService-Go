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

func (r RepositoryContext) BeginTransaction() (*gorm.DB, error) {
	db, err := database.OpenDB()
	if err != nil {
		return nil, err
	}
	return db.Begin(), nil
}

func (r RepositoryContext) Create(tx *gorm.DB, data model.NFAD, status chan error) {

	res := tx.Create(&data)

	if res.RowsAffected == 0 {
		err := errors.New("failed to create")
		status <- err
		close(status)
		return
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) Delete(tx *gorm.DB, id string, status chan error) {

	var nfad model.NFAD
	res := tx.Where("id = ?", id).Delete(&nfad)

	if res.Error != nil {
		status <- res.Error
		close(status)
		return
	}

	if res.RowsAffected == 0 {
		status <- errors.New("not found")
		close(status)
		return
	}

	status <- nil
	close(status)
}

func (r RepositoryContext) Update(tx *gorm.DB, data model.NFAD, status chan error) {

	updateData := map[string]interface{}{
		"PostCode":  data.PostCode,
		"Type":      data.Type,
		"DateStart": data.DateStart,
		"PrevID":    data.PrevID,
		"NextID":    data.NextID,
		"Value":     data.Value,
	}

	res := tx.Model(&model.NFAD{}).Where("id = ?", data.ID).Updates(updateData)

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

func (r RepositoryContext) GetById(tx *gorm.DB, id string, status chan error, data chan *model.NFAD) {
	var nfad model.NFAD
	res := tx.First(&nfad, "id = ?", id)

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

func (r RepositoryContext) GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int, status chan error, data chan []*model.NFAD, totalPages chan int) {
	defer close(status)
	defer close(data)
	defer close(totalPages)

	var nfads []*model.NFAD
	var totalRecords int64

	if err := tx.Model(&model.NFAD{}).Where("post_code = ? AND type = ?", postCode, typeNfad).Count(&totalRecords).Error; err != nil {
		status <- err
		return
	}

	maxPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))
	totalPages <- maxPages

	if err := tx.Where("post_code = ? AND type = ?", postCode, typeNfad).Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&nfads).Error; err != nil {
		status <- err
		return
	}

	if len(nfads) == 0 {
		status <- errors.New("slice is empty")
		return
	}

	data <- nfads
	status <- nil
}

func (r RepositoryContext) GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, status chan error, data chan *model.NFAD) {
	var nfad model.NFAD
	res := tx.Where("post_code = ? AND type = ? AND next_id = ?", postCode, typeNfad, "").First(&nfad)

	if res.Error != nil {
		status <- res.Error
		close(status)
		close(data)
		return
	}

	if res.RowsAffected == 0 {
		data <- nil
		close(data)
		status <- nil
		close(status)
		return
	}

	data <- &nfad
	close(data)
	status <- nil
	close(status)
}

func (r RepositoryContext) GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time, status chan error, data chan []*model.NFAD) {
	date = date.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	for _, typeNfad := range []byte{1, 2, 3, 4} {
		var nfad model.NFAD
		res := tx.Where("post_code = ? AND type = ? AND date_start <= ?", postCode, typeNfad, date).Order("date_start desc").First(&nfad)

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

func (r RepositoryContext) GetByDateRange(tx *gorm.DB, postCode int, startDate, endDate time.Time, status chan error, data chan []*model.NFAD) {
	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	for _, typeNfad := range []byte{1, 2, 3, 4} {
		var nfad model.NFAD
		res := tx.Where("post_code = ? AND type = ? AND date_start <= ?", postCode, typeNfad, endDate).Order("date_start desc").First(&nfad)

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
