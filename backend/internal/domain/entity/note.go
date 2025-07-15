package entity

import (
	"time"
)

type Note struct {
	ID               int64
	UUID             string
	CompanyID        int64
	PersonID         int64   // Sobre quem é a anotação
	UserID           int64   // Manager que escreveu
	Type             string  // "one_on_one", "feedback", "observation"
	Content          string  // Conteúdo com tokens {{person:uuid|nome}}
	FeedbackType     *string // "positive", "constructive", "neutral" - apenas para type="feedback"
	FeedbackCategory *string // "performance", "behavior", "skill", "collaboration" - apenas para type="feedback"
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// IsOneOnOne returns true if the note is a 1:1 meeting record
func (n *Note) IsOneOnOne() bool {
	return n.Type == "one_on_one"
}

// IsFeedback returns true if the note is feedback
func (n *Note) IsFeedback() bool {
	return n.Type == "feedback"
}

// IsObservation returns true if the note is a general observation
func (n *Note) IsObservation() bool {
	return n.Type == "observation"
}

// HasMentions checks if the content contains mention tokens
func (n *Note) HasMentions() bool {
	return len(n.ExtractMentionUUIDs()) > 0
}

// ExtractMentionUUIDs extracts all person UUIDs from mention tokens in content
func (n *Note) ExtractMentionUUIDs() []string {
	// Regex to match {{person:uuid|name}} pattern
	// This is a simple implementation - can be enhanced with proper regex
	var uuids []string
	content := n.Content

	for {
		start := findSubstring(content, "{{person:")
		if start == -1 {
			break
		}

		uuidStart := start + 9 // len("{{person:")
		pipePos := findSubstring(content[uuidStart:], "|")
		if pipePos == -1 {
			break
		}

		uuid := content[uuidStart : uuidStart+pipePos]
		uuids = append(uuids, uuid)

		// Move past this token
		endPos := findSubstring(content[start:], "}}")
		if endPos == -1 {
			break
		}
		content = content[start+endPos+2:]
	}

	return uuids
}

// Helper function to find substring position
func findSubstring(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
