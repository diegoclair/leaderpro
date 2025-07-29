package service

import (
	"context"
	"sync"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type dashboardService struct {
	dm        contract.DataManager
	log       logger.Logger
	authApp   contract.AuthApp
	personApp contract.PersonApp
}

func newDashboardService(infra domain.Infrastructure, authApp contract.AuthApp, personApp contract.PersonApp) contract.DashboardApp {
	return &dashboardService{
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		authApp:   authApp,
		personApp: personApp,
	}
}

func (s *dashboardService) GetDashboardData(ctx context.Context, companyUUID string) (dashboard entity.Dashboard, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get company by UUID to get the ID
	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		return entity.Dashboard{}, resterrors.NewNotFoundError("company not found")
	}

	var (
		wg                                                                  sync.WaitGroup
		peopleErr, totalPeopleErr, oneOnOnesErr, avgFreqErr, lastMeetingErr error
	)

	// Execute all operations in parallel using goroutines
	wg.Add(5)

	// Get people data
	go func() {
		defer wg.Done()
		people, err := s.personApp.GetCompanyPeople(ctx)
		if err != nil {
			peopleErr = err
			return
		}
		dashboard.People = people
	}()

	// Get total people count
	go func() {
		defer wg.Done()
		count, err := s.dm.Person().GetPeopleCountByCompany(ctx, company.ID)
		if err != nil {
			totalPeopleErr = err
			return
		}
		dashboard.Stats.TotalPeople = count
	}()

	// Get one-on-ones this month
	go func() {
		defer wg.Done()
		count, err := s.dm.Note().GetOneOnOnesCountThisMonth(ctx, company.ID)
		if err != nil {
			oneOnOnesErr = err
			return
		}
		dashboard.Stats.OneOnOnesThisMonth = count
	}()

	// Get average frequency
	go func() {
		defer wg.Done()
		avgDays, err := s.dm.Note().GetAverageFrequencyDays(ctx, company.ID)
		if err != nil {
			avgFreqErr = err
			return
		}
		dashboard.Stats.AverageFrequency = avgDays
	}()

	// Get last meeting date
	go func() {
		defer wg.Done()
		lastDate, err := s.dm.Note().GetLastMeetingDate(ctx, company.ID)
		if err != nil {
			lastMeetingErr = err
			return
		}
		dashboard.Stats.LastMeetingDate = lastDate
	}()

	wg.Wait()

	// Check for errors (only fail on critical ones, log others)
	if peopleErr != nil {
		s.log.Errorw(ctx, "error getting company people", logger.Err(peopleErr))
		return dashboard, peopleErr
	}

	// For stats errors, log but don't fail the whole request
	if totalPeopleErr != nil {
		s.log.Errorw(ctx, "error getting people count", logger.Err(totalPeopleErr))
		dashboard.Stats.TotalPeople = 0
	}

	if oneOnOnesErr != nil {
		s.log.Errorw(ctx, "error getting one-on-ones count", logger.Err(oneOnOnesErr))
		dashboard.Stats.OneOnOnesThisMonth = 0
	}

	if avgFreqErr != nil {
		s.log.Errorw(ctx, "error getting average frequency", logger.Err(avgFreqErr))
		dashboard.Stats.AverageFrequency = 0.0
	}

	if lastMeetingErr != nil {
		s.log.Errorw(ctx, "error getting last meeting date", logger.Err(lastMeetingErr))
		dashboard.Stats.LastMeetingDate = nil
	}

	s.log.Infow(ctx, "dashboard data retrieved successfully",
		logger.String("company_uuid", companyUUID),
		logger.Int("people_count", len(dashboard.People)),
		logger.Int64("total_people", dashboard.Stats.TotalPeople),
		logger.Int64("one_on_ones_this_month", dashboard.Stats.OneOnOnesThisMonth),
		logger.Float64("average_frequency_days", dashboard.Stats.AverageFrequency),
	)

	return dashboard, nil
}
