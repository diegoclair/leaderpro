package entity

import (
	"time"
)

// DashboardStats represents the statistics for the dashboard
type DashboardStats struct {
	TotalPeople        int64      `json:"total_people"`
	OneOnOnesThisMonth int64      `json:"one_on_ones_this_month"`
	AverageFrequency   float64    `json:"average_frequency_days"` // in days
	LastMeetingDate    *time.Time `json:"last_meeting_date"`
}

// Dashboard represents the complete dashboard data
type Dashboard struct {
	People []Person       `json:"people"`
	Stats  DashboardStats `json:"stats"`
}