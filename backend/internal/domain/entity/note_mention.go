package entity

import (
	"time"
)

type NoteMention struct {
	ID                int64
	UUID              string
	NoteID            int64
	MentionedPersonID int64  // Pessoa que foi mencionada
	SourcePersonID    int64  // Pessoa sobre quem estava falando (person_id da nota original)
	FullContent       string // Conteúdo completo da nota com tokens
	CreatedAt         time.Time
}

// NoteMentionWithDetails combines mention with related note and person details
type NoteMentionWithDetails struct {
	NoteMention
	Note           Note   // Nota original onde foi feita a menção
	MentionedBy    string // Nome do manager que fez a menção
	SourcePerson   string // Nome da pessoa sobre quem estava falando
	MentionedPerson string // Nome da pessoa que foi mencionada
}

// TimelineEntry represents a unified entry for a person's timeline
type TimelineEntry struct {
	UUID        string    `json:"uuid"`
	Type        string    `json:"type"`        // "one_on_one", "feedback", "observation"
	Content     string    `json:"content"`
	AuthorName  string    `json:"author_name"`
	CreatedAt   time.Time `json:"created_at"`
	
	// For feedback notes
	FeedbackType     *string `json:"feedback_type,omitempty"`
	FeedbackCategory *string `json:"feedback_category,omitempty"`
	
	// For mentions
	SourcePersonName *string `json:"source_person_name,omitempty"` // Pessoa sobre quem falou
}

// MentionEntry represents notes where a person was mentioned (feedbacks received)
type MentionEntry struct {
	UUID             string    `json:"uuid"`
	Type             string    `json:"type"`             // "one_on_one", "feedback", "observation"
	Content          string    `json:"content"`
	FeedbackType     *string   `json:"feedback_type,omitempty"`
	FeedbackCategory *string   `json:"feedback_category,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	PersonID         string    `json:"person_id"`   // UUID da pessoa sobre quem a nota foi feita
	PersonName       string    `json:"person_name"` // Nome da pessoa sobre quem a nota foi feita
	Mentions         []struct {
		ID          string `json:"id"`
		PersonID    string `json:"person_id"`
		PersonName  string `json:"person_name"`
		StartIndex  int    `json:"start_index"`
		EndIndex    int    `json:"end_index"`
	} `json:"mentions,omitempty"`
}