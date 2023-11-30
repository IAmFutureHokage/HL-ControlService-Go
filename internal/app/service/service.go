package service

import (
	"context"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/domain/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Repository interface {
	AddControlValue(ctx context.Context, value model.ControlValue) error
	RemoveControlValue(ctx context.Context, id string) error
	UpdateControlValue(ctx context.Context, values []model.ControlValue) error
	// GetById(tx *gorm.DB, id string) (*model.NFAD, error)
	// GetByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte, pageNumber, pageSize int) (int, []*model.NFAD, error)
	// GetActiveByPostCodeAndType(tx *gorm.DB, postCode int, typeNfad byte) (*model.NFAD, error)
	// GetByPostCodeAndDate(tx *gorm.DB, postCode int, date time.Time) ([]*model.NFAD, error)
	// GetByDateRange(tx *gorm.DB, postCode int, startDate, endDate time.Time) ([]*model.NFAD, error)
}

type HydrologyStatsService struct {
	repo *repository.HydrologyStatsRepository
	pb.UnimplementedHydrologyStatsServiceServer
}

func NewHydrologyStatsService(repo *repository.HydrologyStatsRepository) *HydrologyStatsService {
	return &HydrologyStatsService{repo: repo}
}

func (s *HydrologyStatsService) AddControlValue(ctx context.Context, req *pb.AddControlValueRequest) (*pb.AddControlValueResponse, error) {

	controlValue := model.ControlValue{
		ID:        uuid.New().String(),
		PostCode:  req.GetPostCode(),
		Type:      model.ControlValueType(req.GetType()),
		DateStart: req.GetDateStart().AsTime().Truncate(24 * time.Hour),
		Value:     req.GetValue(),
	}

	if err := s.repo.AddControlValue(ctx, controlValue); err != nil {
		return nil, err
	}

	return &pb.AddControlValueResponse{
		ControlValue: &pb.ControlValue{
			Id:        controlValue.ID,
			PostCode:  controlValue.PostCode,
			Type:      pb.ControlValueType(controlValue.Type),
			DateStart: timestamppb.New(controlValue.DateStart),
			Value:     controlValue.Value,
		},
	}, nil
}

func (s *HydrologyStatsService) RemoveControlValue(ctx context.Context, req *pb.RemoveControlValueRequest) (*pb.RemoveControlValueResponse, error) {

	err := s.repo.RemoveControlValue(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.RemoveControlValueResponse{Success: true}, nil
}

func (s *HydrologyStatsService) UpdateControlValue(ctx context.Context, req *pb.UpdateControlValueRequest) (*pb.UpdateControlValueResponse, error) {

	var controlValues []model.ControlValue

	for _, cv := range req.GetControlValues() {
		controlValues = append(controlValues, model.ControlValue{
			ID:        cv.GetId(),
			PostCode:  cv.GetPostCode(),
			Type:      model.ControlValueType(cv.GetType()),
			DateStart: cv.DateStart.AsTime().Truncate(24 * time.Hour),
			Value:     cv.GetValue(),
		})
	}

	err := s.repo.UpdateControlValues(ctx, controlValues)
	if err != nil {
		return nil, err
	}

	var updatedControlValues []*pb.ControlValue
	for _, val := range controlValues {
		updatedControlValues = append(updatedControlValues, &pb.ControlValue{
			Id:        val.ID,
			PostCode:  val.PostCode,
			Type:      pb.ControlValueType(val.Type),
			DateStart: timestamppb.New(val.DateStart),
			Value:     val.Value,
		})
	}

	return &pb.UpdateControlValueResponse{ControlValues: updatedControlValues}, nil
}
