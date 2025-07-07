package entity

import (
	"time"
)

type Feedback struct {
	ID            int64
	UUID          string
	CompanyID     int64
	PersonID      int64      // Who receives the feedback
	GivenBy       int64      // Who gave the feedback
	OneOnOneID    *int64     // If feedback was given during a 1:1
	
	Type          string     // positive, constructive, observation
	Category      string     // performance, behavior, skill, collaboration
	Content       string
	
	// For @mentions and cross-references
	MentionedFrom string     // e.g., "1:1 with Maria"
	MentionedDate time.Time
	
	// Metadata
	CreatedAt     time.Time
	UpdatedAt     time.Time
	IsPrivate     bool       // If true, only visible to manager
}

// FeedbackSummary aggregates feedback for a person over a period
type FeedbackSummary struct {
	PersonID      int64
	Period        string     // e.g., "Q1 2024"
	TotalCount    int
	PositiveCount int
	Constructive  int
	TopCategories []string
	Highlights    []string   // Key feedback excerpts
}