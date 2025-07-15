package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

// DashboardStatsResponse represents the statistics for the dashboard
type DashboardStatsResponse struct {
	TotalPeople        int64      `json:"total_people"`
	OneOnOnesThisMonth int64      `json:"one_on_ones_this_month"`
	AverageFrequency   float64    `json:"average_frequency_days"` // in days
	LastMeetingDate    *time.Time `json:"last_meeting_date"`
}

// DashboardResponse represents the complete dashboard data
type DashboardResponse struct {
	People []PersonResponse        `json:"people"`
	Stats  DashboardStatsResponse  `json:"stats"`
}

// FillFromEntity fills the dashboard response from entity
func (r *DashboardResponse) FillFromEntity(dashboard entity.Dashboard) {
	// Fill people array
	r.People = make([]PersonResponse, len(dashboard.People))
	for i, person := range dashboard.People {
		r.People[i].FillFromEntity(person)
	}

	// Fill stats
	r.Stats = DashboardStatsResponse{
		TotalPeople:        dashboard.Stats.TotalPeople,
		OneOnOnesThisMonth: dashboard.Stats.OneOnOnesThisMonth,
		AverageFrequency:   dashboard.Stats.AverageFrequency,
		LastMeetingDate:    dashboard.Stats.LastMeetingDate,
	}
}