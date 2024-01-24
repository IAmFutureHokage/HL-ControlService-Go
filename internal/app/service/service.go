package service

import (
	"context"

	//"fmt"
	"math"
	"sync"
	"time"

	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/model"
	"github.com/IAmFutureHokage/HL-ControlService-Go/internal/app/repository"
	pb "github.com/IAmFutureHokage/HL-ControlService-Go/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type HydrologyStatsRepository interface {
	AddControlValue(ctx context.Context, value model.ControlValue) error
	RemoveControlValue(ctx context.Context, id string) error
	UpdateControlValue(ctx context.Context, values []model.ControlValue) error
	GetControlValues(ctx context.Context, postCode string, controlType model.ControlValueType, page, pageSize int) ([]model.ControlValue, int, error)
	GetControlValuesByDay(ctx context.Context, postCode string, date time.Time) ([]model.ControlValue, error)
	GetControlValuesByDateInterval(ctx context.Context, postCode string, dateStart, dateEnd time.Time) ([]model.ControlValue, error)
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

func (s *HydrologyStatsService) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {

	startDate := req.GetStartDate().AsTime().Truncate(24 * time.Hour)
	endDate := req.GetEndDate().AsTime().Truncate(24 * time.Hour)
	postCode := req.GetPostCode()
	graphPoints := req.GetGraphPoints()

	controlValuesCh := make(chan []model.ControlValue)
	waterLevelsCh := make(chan []model.Waterlevel)
	startWaterIntervalCh := make(chan time.Time, 1)
	errorCh := make(chan error)
	doneCh := make(chan struct{})
	once := sync.Once{}

	go func() {
		defer func() {
			close(controlValuesCh)
			once.Do(func() { close(doneCh) })
		}()

		controlValues, err := s.repo.GetControlValuesByDateInterval(ctx, postCode, startDate, endDate)
		if err != nil {
			errorCh <- err
			return
		}
		controlValuesCh <- controlValues
	}()

	go func() {
		defer func() {
			close(waterLevelsCh)
			once.Do(func() { close(doneCh) })
		}()

		waterLevels, err := s.repo.GetWaterlevelsByDateInterval(ctx, postCode, startDate, endDate)
		if err != nil {
			errorCh <- err
			return
		}
		waterLevelsCh <- waterLevels
	}()

	go func() {
		defer func() {
			close(startWaterIntervalCh)
			once.Do(func() { close(doneCh) })
		}()

		startWaterInterval, err := s.repo.GetStartInterval(ctx, postCode)
		if err != nil {
			errorCh <- err
			return
		}
		startWaterIntervalCh <- startWaterInterval
	}()

	<-doneCh

	controlValues := <-controlValuesCh
	waterLevels := <-waterLevelsCh
	startWaterInterval := <-startWaterIntervalCh
	startInterval := timestamppb.New(startWaterInterval)

	var funcCheck bool
	if uint32(len(waterLevels)) < graphPoints {
		graphPoints = uint32(len(waterLevels))
		funcCheck = true
	}

	allStats := make([]*pb.StatsDay, graphPoints)
	var dates []time.Time

	if funcCheck {
		allStats, dates = allWaterPoints(waterLevels, allStats, graphPoints)
	} else {
		allStats, dates, graphPoints = rangedWaterPoints(waterLevels, allStats, graphPoints)
	}

	var norm []model.ControlValue
	var floodplain []model.ControlValue
	var adverse []model.ControlValue
	var dangerous []model.ControlValue

	for _, cv := range controlValues {
		cvCopy := cv
		if cvCopy.Type == 1 {
			norm = append(norm, cvCopy)
		}
		if cvCopy.Type == 2 {
			floodplain = append(floodplain, cvCopy)
		}
		if cvCopy.Type == 3 {
			adverse = append(adverse, cvCopy)
		}
		if cvCopy.Type == 4 {
			dangerous = append(dangerous, cvCopy)
		}
	}

	for i := 0; i < int(len(allStats)); i++ {
		allStats[i].Norm = getTypeLevel(norm, dates[i])
		allStats[i].Floodplain = getTypeLevel(floodplain, dates[i])
		allStats[i].Adverse = getTypeLevel(adverse, dates[i])
		allStats[i].Dangerous = getTypeLevel(dangerous, dates[i])
	}

	return &pb.GetStatsResponse{
		StartInterval: startInterval,
		Stats:         allStats,
	}, nil
}

func getTypeLevel(values []model.ControlValue, date time.Time) uint32 {
	if len(values) == 0 {
		return 0
	}
	if len(values) == 1 {
		return uint32(values[0].Value)
	}
	for i := range values {
		if values[i].DateStart.Before(date) && values[i+1].DateStart.After(date) || values[i].DateStart == date {
			return uint32(values[i].Value)
		}
		if i == len(values)-2 {
			return uint32(values[i+1].Value)
		}
	}

	return 0
}

func allWaterPoints(waterLevels []model.Waterlevel, allStats []*pb.StatsDay, points uint32) ([]*pb.StatsDay, []time.Time) {
	dates := make([]time.Time, points)
	for i := 0; i < int(points); i++ {
		var stat pb.StatsDay
		stat.Date = timestamppb.New(waterLevels[i].Date)
		stat.Waterlevel = wrapperspb.Int32(int32(waterLevels[i].Waterlevel))
		allStats[i] = &stat
		dates[i] = waterLevels[i].Date
	}

	return allStats, dates
}

func rangedWaterPoints(waterLevels []model.Waterlevel, allStats []*pb.StatsDay, points uint32) ([]*pb.StatsDay, []time.Time, uint32) {
	step := float64(len(waterLevels)) / float64(points)
	dates := make([]time.Time, points)
	h := 0

	for i := float64(0); int(math.Round(i)) < len(waterLevels); i += step {
		var stat pb.StatsDay
		stat.Date = timestamppb.New(waterLevels[int(math.Round(i))].Date)
		stat.Waterlevel = wrapperspb.Int32(int32(waterLevels[int(math.Round(i))].Waterlevel))
		allStats[h] = &stat
		dates[h] = waterLevels[int(math.Round(i))].Date
		h++
		if h == int(points) {
			return allStats, dates, points
		}
	}

	points = uint32(len(allStats))

	if points != 100 {
		newAllStats := make([]*pb.StatsDay, points)
		copy(newAllStats, allStats)
		allStats = newAllStats

		newDates := make([]time.Time, points)
		copy(newDates, dates)
		dates = newDates
	}

	return allStats, dates, points
}
