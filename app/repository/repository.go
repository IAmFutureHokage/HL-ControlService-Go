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
	err := tx.Create(&data).Error
	status <- err

	close(status)
}

func (r RepositoryContext) Delete(tx *gorm.DB, id string, status chan error) {

	var nfad model.NFAD
	err := tx.Where("id = ?", id).Delete(&nfad).Error
	status <- err

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

	err := tx.Model(&model.NFAD{}).Where("id = ?", data.ID).Updates(updateData).Error
	status <- err

	close(status)
}

func (r RepositoryContext) GetById(tx *gorm.DB, id string, status chan error, data chan *model.NFAD) {

	var nfad model.NFAD
	err := tx.First(&nfad, "id = ?", id).Error

	status <- err
	close(status)

	if err == nil {
		data <- &nfad
	}
	close(data)
}

func (r RepositoryContext) GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, status chan error, data chan *model.NFAD) {

	var nfad model.NFAD

	err := tx.Where("post_code = ? AND type = ? AND (next_id IS NULL OR next_id = '')", postCode, typeNfad).First(&nfad).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	status <- err
	close(status)

	if err == nil {
		data <- &nfad
	}
	close(data)
}

func (r RepositoryContext) GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int, status chan error, data chan []*model.NFAD, totalPages chan int) {

	var nfads []*model.NFAD
	var totalRecords int64

	err := tx.Model(&model.NFAD{}).Where("post_code = ? AND type = ?", postCode, typeNfad).Count(&totalRecords).Error
	if err == nil {
		maxPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))
		totalPages <- maxPages

		offset := (pageNumber - 1) * pageSize
		err = tx.Where("post_code = ? AND type = ?", postCode, typeNfad).Offset(offset).Limit(pageSize).Find(&nfads).Error
		if err == nil {
			data <- nfads
		}
	}

	status <- err
	close(status)
	close(totalPages)
	close(data)
}

func (r RepositoryContext) GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time, status chan error, data chan []*model.NFAD) {

	date = date.Truncate(24 * time.Hour)

	var nfads []*model.NFAD
	var err error

	for _, typeNfad := range []byte{1, 2, 3, 4} {
		var nfad model.NFAD
		err := tx.Where("post_code = ? AND type = ? AND date_start <= ?", postCode, typeNfad, date).Order("date_start desc").First(&nfad).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			break
		}

		nfads = append(nfads, &nfad)
	}

	if err == nil && len(nfads) == 0 {
		err = errors.New("no records found for the specified post code and date")
	}

	status <- err
	close(status)

	if err == nil {
		data <- nfads
	}
	close(data)
}

func (r RepositoryContext) GetByDateRange(tx *gorm.DB, postCode int, startDate, endDate time.Time, status chan error, data chan []*model.NFAD) {

	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	err := tx.Where("post_code = ? AND date_start <= ? AND (next_id IS NULL OR next_id = '' OR next_id IN (SELECT id FROM nfads WHERE date_start >= ?))", postCode, endDate, startDate).Order("date_start desc").Find(&nfads).Error

	if err == nil && len(nfads) == 0 {
		err = errors.New("no records found for the specified range")
	}

	status <- err
	close(status)

	if err == nil {
		data <- nfads
	}
	close(data)
}
