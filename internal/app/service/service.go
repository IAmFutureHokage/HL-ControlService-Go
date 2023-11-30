package service

import (
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/internal/proto"
)

// type Repository interface {
// 	Create(tx *gorm.DB, data model.NFAD) error
// 	Delete(tx *gorm.DB, id string) error
// 	Update(tx *gorm.DB, data model.NFAD) error
// 	GetById(tx *gorm.DB, id string) (*model.NFAD, error)
// 	GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int) (int, []*model.NFAD, error)
// 	GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte) (*model.NFAD, error)
// 	GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time) ([]*model.NFAD, error)
// 	GetByDateRange(tx *gorm.DB, postCode int, startDate, endDate time.Time) ([]*model.NFAD, error)
// }

type HydrologyStatsService struct {
	pb.UnimplementedHydrologyStatsServiceServer
}

// func (*ServerContext) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
// 	repo := new(repository.RepositoryContext)

// 	tx, err := repo.BeginTransaction()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(req.Data) < 1 {
// 		return nil, errors.New("no data provided")
// 	}

// 	if len(req.Data) != 1 {

// 		for i := 0; i < len(req.Data)-1; i++ {
// 			current := req.Data[i]
// 			next := req.Data[i+1]

// 			if current.PostCode != next.PostCode ||
// 				current.Type != next.Type ||
// 				current.Id != next.PrevId ||
// 				!next.DateStart.AsTime().Truncate(24*time.Hour).After(current.DateStart.AsTime().Truncate(24*time.Hour)) {
// 				return nil, errors.New("bad data")
// 			}
// 		}
// 	}

// 	status := make(chan error)
// 	close(status)

// 	for _, pbNFAD := range req.Data {

// 		nfad := model.NFAD{
// 			ID:        pbNFAD.Id,
// 			PostCode:  pbNFAD.PostCode,
// 			Type:      model.ControlType(pbNFAD.Type),
// 			DateStart: pbNFAD.DateStart.AsTime().Truncate(24 * time.Hour),
// 			Value:     pbNFAD.Value,
// 			PrevID:    pbNFAD.PrevId,
// 			NextID:    pbNFAD.NextId,
// 		}

// 		status = make(chan error)
// 		go repo.Update(tx, nfad, status)
// 		err = <-status
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}
// 	}

// 	tx.Commit()

// 	return &pb.UpdateResponse{
// 		Data: req.Data,
// 	}, nil
// }
