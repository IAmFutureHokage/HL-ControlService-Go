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

func (r RepositoryContext) Create(tx *gorm.DB, data model.NFAD) error {
	return tx.Create(&data).Error
}

func (r RepositoryContext) Delete(tx *gorm.DB, id string) error {
	return tx.Where("id = ?", id).Delete(&model.NFAD{}).Error
}

func (r RepositoryContext) Update(tx *gorm.DB, data model.NFAD) error {
	updateData := map[string]interface{}{
		"PostCode":  data.PostCode,
		"Type":      data.Type,
		"DateStart": data.DateStart,
		"PrevID":    data.PrevID,
		"NextID":    data.NextID,
		"Value":     data.Value,
	}

	return tx.Model(&model.NFAD{}).Where("id = ?", data.ID).Updates(updateData).Error
}

func (r RepositoryContext) GetById(tx *gorm.DB, id string) (*model.NFAD, error) {
	var nfad model.NFAD
	err := tx.First(&nfad, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return &nfad, nil
}

func (r RepositoryContext) GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte) (*model.NFAD, error) {
	var nfad model.NFAD

	err := tx.Where("post_code = ? AND type = ? AND (next_id IS NULL OR next_id = '')", postCode, typeNfad).First(&nfad).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &nfad, nil
}

func (r RepositoryContext) GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int) (int, []*model.NFAD, error) {

	var nfads []*model.NFAD
	var totalRecords int64

	err := tx.Model(&model.NFAD{}).Where("post_code = ? AND type = ?", postCode, typeNfad).Count(&totalRecords).Error

	if err != nil {
		return 0, nil, err
	}

	maxPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))
	totalPages := maxPages

	offset := (pageNumber - 1) * pageSize
	err = tx.Where("post_code = ? AND type = ?", postCode, typeNfad).Offset(offset).Limit(pageSize).Find(&nfads).Error

	if err != nil {
		return 0, nil, err
	}
	return totalPages, nfads, nil
}

func (r RepositoryContext) GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time) ([]*model.NFAD, error) {
	date = date.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	err := tx.Raw(`
        SELECT DISTINCT ON (type) *
        FROM nfads
        WHERE post_code = ? AND date_start <= ?
        ORDER BY type, date_start DESC
    `, postCode, date).Scan(&nfads).Error

	if err != nil {
		return nil, err
	}

	if len(nfads) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return nfads, nil
}

func (r RepositoryContext) GetByDateRange(tx *gorm.DB, postCode int, startDate, endDate time.Time) ([]*model.NFAD, error) {

	startDate = startDate.Truncate(24 * time.Hour)
	endDate = endDate.Truncate(24 * time.Hour)

	var nfads []*model.NFAD

	err := tx.Where("post_code = ? AND date_start <= ? AND (next_id IS NULL OR next_id = '' OR next_id IN (SELECT id FROM nfads WHERE date_start >= ?))", postCode, endDate, startDate).Order("date_start desc").Find(&nfads).Error

	if err != nil {
		return nil, err
	}
	return nfads, nil
}
