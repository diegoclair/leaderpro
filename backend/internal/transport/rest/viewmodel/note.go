package viewmodel

import (
	"time"

	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type CreateNoteRequest struct {
	Type             string            `json:"type" validate:"required,oneof=one_on_one feedback observation"`
	Content          string            `json:"content" validate:"required,min=1"`
	FeedbackType     *string           `json:"feedback_type,omitempty" validate:"omitempty,oneof=positive constructive neutral"`
	FeedbackCategory *string           `json:"feedback_category,omitempty" validate:"omitempty,oneof=performance behavior skill collaboration"`
	MentionedPeople  []MentionedPerson `json:"mentioned_people,omitempty"`
}

type MentionedPerson struct {
	UUID string `json:"uuid" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type UpdateNoteRequest struct {
	Type             string            `json:"type" validate:"required,oneof=one_on_one feedback observation"`
	Content          string            `json:"content" validate:"required,min=1"`
	FeedbackType     *string           `json:"feedback_type,omitempty" validate:"omitempty,oneof=positive constructive neutral"`
	FeedbackCategory *string           `json:"feedback_category,omitempty" validate:"omitempty,oneof=performance behavior skill collaboration"`
	MentionedPeople  []MentionedPerson `json:"mentioned_people,omitempty"`
}

type NoteResponse struct {
	UUID             string    `json:"uuid"`
	Type             string    `json:"type"`
	Content          string    `json:"content"`
	FeedbackType     *string   `json:"feedback_type,omitempty"`
	FeedbackCategory *string   `json:"feedback_category,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TimelineResponse struct {
	UUID             string    `json:"uuid"`
	Type             string    `json:"type"`        // "one_on_one", "feedback", "observation"
	Content          string    `json:"content"`
	AuthorName       string    `json:"author_name"`
	CreatedAt        time.Time `json:"created_at"`
	FeedbackType     *string   `json:"feedback_type,omitempty"`
	FeedbackCategory *string   `json:"feedback_category,omitempty"`
	SourcePersonName *string   `json:"source_person_name,omitempty"` // Para mentions: pessoa sobre quem falou
}

type MentionResponse struct {
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

func (r *CreateNoteRequest) ToEntity() entity.Note {
	return entity.Note{
		Type:             r.Type,
		Content:          r.Content,
		FeedbackType:     r.FeedbackType,
		FeedbackCategory: r.FeedbackCategory,
	}
}

func (r *UpdateNoteRequest) ToEntity() entity.Note {
	return entity.Note{
		Type:             r.Type,
		Content:          r.Content,
		FeedbackType:     r.FeedbackType,
		FeedbackCategory: r.FeedbackCategory,
	}
}

func (r *NoteResponse) FillFromEntity(note entity.Note) {
	r.UUID = note.UUID
	r.Type = note.Type
	r.Content = note.Content
	r.FeedbackType = note.FeedbackType
	r.FeedbackCategory = note.FeedbackCategory
	r.CreatedAt = note.CreatedAt
	r.UpdatedAt = note.UpdatedAt
}

func (r *TimelineResponse) FillFromTimelineEntry(entry entity.TimelineEntry) {
	r.UUID = entry.UUID
	r.Type = entry.Type
	r.Content = entry.Content
	r.AuthorName = entry.AuthorName
	r.CreatedAt = entry.CreatedAt
	r.FeedbackType = entry.FeedbackType
	r.FeedbackCategory = entry.FeedbackCategory
	r.SourcePersonName = entry.SourcePersonName
}

func (r *MentionResponse) FillFromMentionEntry(entry entity.MentionEntry) {
	r.UUID = entry.UUID
	r.Type = entry.Type
	r.Content = entry.Content
	r.FeedbackType = entry.FeedbackType
	r.FeedbackCategory = entry.FeedbackCategory
	r.CreatedAt = entry.CreatedAt
	r.PersonID = entry.PersonID
	r.PersonName = entry.PersonName
	
	// Convert mentions
	r.Mentions = make([]struct {
		ID          string `json:"id"`
		PersonID    string `json:"person_id"`
		PersonName  string `json:"person_name"`
		StartIndex  int    `json:"start_index"`
		EndIndex    int    `json:"end_index"`
	}, len(entry.Mentions))
	
	for i, mention := range entry.Mentions {
		r.Mentions[i].ID = mention.ID
		r.Mentions[i].PersonID = mention.PersonID
		r.Mentions[i].PersonName = mention.PersonName
		r.Mentions[i].StartIndex = mention.StartIndex
		r.Mentions[i].EndIndex = mention.EndIndex
	}
}