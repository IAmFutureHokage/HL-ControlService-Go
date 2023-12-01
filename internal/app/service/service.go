package service

import (
	"context"
	"math"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type HydrologyStatsRepository interface {
	AddControlValue(ctx context.Context, value model.ControlValue) error
	RemoveControlValue(ctx context.Context, id string) error
	UpdateControlValue(ctx context.Context, values []model.ControlValue) error
	GetControlValues(ctx context.Context, postCode string, controlType model.ControlValueType, page, pageSize int) ([]model.ControlValue, int, error)
	GetControlValuesByDay(ctx context.Context, postCode string, date time.Time) ([]model.ControlValue, error)
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

func (s *HydrologyStatsService) GetControlValues(ctx context.Context, req *pb.GetControlValuesRequest) (*pb.GetControlValuesResponse, error) {

	pageSize := 50

	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}

	controlValues, totalCount, err := s.repo.GetControlValues(ctx, req.GetPostCode(), model.ControlValueType(req.GetType()), page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbControlValues []*pb.ControlValue

	for _, cv := range controlValues {
		pbControlValues = append(pbControlValues, &pb.ControlValue{
			Id:        cv.ID,
			PostCode:  cv.PostCode,
			Type:      pb.ControlValueType(cv.Type),
			DateStart: timestamppb.New(cv.DateStart),
			Value:     cv.Value,
		})
	}

	maxPage := uint32(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &pb.GetControlValuesResponse{
		Page:          uint32(page),
		MaxPage:       maxPage,
		ControlValues: pbControlValues,
	}, nil
}

func (s *HydrologyStatsService) CheckWaterLevel(ctx context.Context, req *pb.CheckWaterLevelRequest) (*pb.CheckWaterLevelResponse, error) {

	date := req.GetDate().AsTime().Truncate(24 * time.Hour)

	controlValues, err := s.repo.GetControlValuesByDay(ctx, req.GetPostCode(), date)
	if err != nil {
		return nil, err
	}

	var excessType uint32

	for _, cv := range controlValues {
		if cv.Value < req.GetValue() {
			excessType = uint32(cv.Type)
			break
		}
	}

	return &pb.CheckWaterLevelResponse{
		Excess: excessType,
	}, nil
}

func (s *HydrologyStatsService) GetStatByDay(ctx context.Context, req *pb.GetStatByDayRequest) (*pb.GetStatByDayResponse, error) {

	date := req.GetDate().AsTime().Truncate(24 * time.Hour)

	controlValues, err := s.repo.GetControlValuesByDay(ctx, req.GetPostCode(), date)
	if err != nil {
		return nil, err
	}

	norm, floodplain, adverse, dangerous := 0, 0, 0, 0

	for _, cv := range controlValues {
		if cv.Type.ToByte() == 1 {
			norm = int(cv.Value)
		}
		if cv.Type.ToByte() == 2 {
			floodplain = int(cv.Value)
		}
		if cv.Type.ToByte() == 3 {
			adverse = int(cv.Value)
		}
		if cv.Type.ToByte() == 4 {
			dangerous = int(cv.Value)
		}
	}

	return &pb.GetStatByDayResponse{
		Stat: &pb.StatByDay{
			Date:       timestamppb.New(date),
			Norm:       uint32(norm),
			Floodplain: uint32(floodplain),
			Adverse:    uint32(adverse),
			Dangerous:  uint32(dangerous),
		},
	}, nil
}
