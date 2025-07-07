package entity

import (
	"time"
)

type OneOnOne struct {
	ID              int64
	UUID            string
	CompanyID       int64
	PersonID        int64
	ManagerID       int64
	ScheduledDate   time.Time
	ActualDate      *time.Time
	Duration        int // in minutes
	Location        string
	Status          string // scheduled, completed, cancelled
	
	// Meeting content
	Agenda          string
	DiscussionNotes string
	ActionItems     string
	PrivateNotes    string // Manager's private notes
	
	// AI suggestions (stored as JSON)
	AIContext       string
	AISuggestions   string
	
	// Metadata
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CompletedAt     *time.Time
}

// IsCompleted returns true if the 1:1 has been completed
func (o *OneOnOne) IsCompleted() bool {
	return o.Status == "completed" && o.CompletedAt != nil
}

// IsOverdue returns true if the scheduled date has passed and the meeting isn't completed
func (o *OneOnOne) IsOverdue() bool {
	return o.ScheduledDate.Before(time.Now()) && !o.IsCompleted()
}