package service

import (
	"context"
	"fmt"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type aiApp struct {
	dm         contract.DataManager
	aiProvider contract.AIProvider
	authApp    contract.AuthApp
	log        logger.Logger
}

func newAIApp(infra domain.Infrastructure, aiProvider contract.AIProvider, authApp contract.AuthApp) contract.AIApp {
	return &aiApp{
		dm:         infra.DataManager(),
		aiProvider: aiProvider,
		authApp:    authApp,
		log:        infra.Logger(),
	}
}

func (s *aiApp) ChatWithLeadershipCoach(ctx context.Context, req entity.ChatRequest) (entity.ChatResponse, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return entity.ChatResponse{}, fmt.Errorf("failed to get logged user: %w", err)
	}

	companyUUID, err := s.authApp.GetCompanyFromContext(ctx)
	if err != nil {
		return entity.ChatResponse{}, fmt.Errorf("failed to get company UUID: %w", err)
	}

	company, err := s.dm.Company().GetCompanyByUUID(ctx, companyUUID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get company", logger.Err(err))
		return entity.ChatResponse{}, fmt.Errorf("failed to get company: %w", err)
	}
	companyID := company.ID

	var personID *int64
	if req.PersonUUID != nil {
		person, err := s.dm.Person().GetPersonByUUID(ctx, *req.PersonUUID)
		if err != nil {
			s.log.Errorw(ctx, "failed to get person", logger.Err(err))
			return entity.ChatResponse{}, fmt.Errorf("failed to get person: %w", err)
		}
		personID = &person.ID
	}

	prompt, err := s.dm.AI().GetActivePromptByType(ctx, domain.AIPromptTypeLeadershipCoach)
	if err != nil {
		s.log.Errorw(ctx, "failed to get leadership coach prompt", logger.Err(err))
		return entity.ChatResponse{}, fmt.Errorf("failed to get leadership coach prompt: %w", err)
	}

	var contextPrompt string
	if personID != nil {
		personContext, err := s.GetPersonContext(ctx, *personID)
		if err != nil {
			// Continue without person context on error
			contextPrompt = ""
		} else {
			contextPrompt = s.buildContextPrompt(personContext)
		}
	}

	start := time.Now()
	response, err := s.aiProvider.Chat(ctx, req, prompt.Prompt, contextPrompt)
	if err != nil {
		return entity.ChatResponse{}, fmt.Errorf("ai provider error: %w", err)
	}
	responseTime := int(time.Since(start).Milliseconds())

	usage := entity.AIUsageTracker{
		UserID:         userID,
		CompanyID:      companyID,
		PromptID:       prompt.ID,
		PersonID:       personID,
		RequestType:    "chat",
		TokensUsed:     response.Usage.TotalTokens,
		CostUSD:        response.Usage.CostUSD,
		ResponseTimeMs: responseTime,
	}

	createdUsage, err := s.dm.AI().CreateUsage(ctx, usage)
	if err != nil {
		// Don't fail request on usage tracking error
		s.log.Errorw(ctx, "Failed to create usage record", logger.Err(err))
	}

	conversation := entity.AIConversation{
		UsageID:     createdUsage.ID,
		UserMessage: req.Message,
		AIResponse:  response.Response,
	}

	_, err = s.dm.AI().CreateConversation(ctx, conversation)
	if err != nil {
		// Don't fail request on conversation save error
		s.log.Errorw(ctx, "Failed to save conversation", logger.Err(err))
	}

	response.UsageID = createdUsage.ID
	return response, nil
}

func (s *aiApp) ExtractAttributesFromNote(ctx context.Context, noteID int64) (entity.AttributesResponse, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	note, err := s.dm.Note().GetNoteByID(ctx, noteID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get note", logger.Err(err))
		return entity.AttributesResponse{}, fmt.Errorf("failed to get note: %w", err)
	}

	prompt, err := s.dm.AI().GetActivePromptByType(ctx, domain.AIPromptTypeAttributeExtraction)
	if err != nil {
		s.log.Errorw(ctx, "failed to get extraction prompt", logger.Err(err))
		return entity.AttributesResponse{}, fmt.Errorf("failed to get extraction prompt: %w", err)
	}

	extractionReq := entity.ExtractionRequest{
		PersonID: note.PersonID,
		NoteID:   noteID,
		Content:  note.Content,
	}

	start := time.Now()
	extractedAttributes, usageInfo, err := s.aiProvider.ExtractAttributes(ctx, extractionReq, prompt.Prompt)
	if err != nil {
		return entity.AttributesResponse{}, fmt.Errorf("ai extraction error: %w", err)
	}
	responseTime := int(time.Since(start).Milliseconds())

	userID, _ := s.authApp.GetLoggedUserID(ctx)

	person, err := s.dm.Person().GetPersonByID(ctx, note.PersonID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get person for usage tracking", logger.Err(err))
		return entity.AttributesResponse{}, fmt.Errorf("failed to get person for usage tracking: %w", err)
	}
	companyID := person.CompanyID

	usage := entity.AIUsageTracker{
		UserID:         userID,
		CompanyID:      companyID,
		PromptID:       prompt.ID,
		PersonID:       &note.PersonID,
		RequestType:    "extraction",
		TokensUsed:     usageInfo.TotalTokens,
		CostUSD:        usageInfo.CostUSD,
		ResponseTimeMs: responseTime,
	}

	createdUsage, err := s.dm.AI().CreateUsage(ctx, usage)
	if err != nil {
		s.log.Errorw(ctx, "Failed to create usage record", logger.Err(err))
	}

	var attributes []entity.PersonAttribute
	if len(extractedAttributes) > 0 {
		err = s.dm.Person().BulkUpsertPersonAttributes(ctx, note.PersonID, extractedAttributes, "ai_extracted", &noteID)
		if err != nil {
			s.log.Errorw(ctx, "failed to save extracted attributes", logger.Err(err))
			return entity.AttributesResponse{}, fmt.Errorf("failed to save extracted attributes: %w", err)
		}

		for key, value := range extractedAttributes {
			attr := entity.PersonAttribute{
				PersonID:            note.PersonID,
				AttributeKey:        key,
				AttributeValue:      value,
				Source:              "ai_extracted",
				ExtractedFromNoteID: &noteID,
			}
			attributes = append(attributes, attr)
		}
	}

	return entity.AttributesResponse{
		Attributes: attributes,
		UsageID:    createdUsage.ID,
		Usage:      usageInfo,
	}, nil
}

func (s *aiApp) GetPersonContext(ctx context.Context, personID int64) (entity.PersonAIContext, error) {
	person, err := s.dm.Person().GetPersonByID(ctx, personID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get person", logger.Err(err))
		return entity.PersonAIContext{}, fmt.Errorf("failed to get person: %w", err)
	}

	attributes, err := s.dm.Person().GetPersonAttributesMap(ctx, personID)
	if err != nil {
		s.log.Errorw(ctx, "failed to get person attributes", logger.Err(err))
		attributes = make(map[string]string)
	}

	notes, err := s.dm.Note().GetNotesByPersonIDPaginated(ctx, personID, 1, 50)
	if err != nil {
		s.log.Errorw(ctx, "failed to get person notes", logger.Err(err))
		notes = []entity.Note{}
	}

	var lastMeeting *entity.PersonLastMeeting
	for _, note := range notes {
		if note.Type == domain.NoteTypeOneOnOne {
			lastMeeting = &entity.PersonLastMeeting{
				Date:  note.CreatedAt,
				Notes: note.Content,
			}
			break
		}
	}

	return entity.PersonAIContext{
		Person:      person,
		Attributes:  attributes,
		RecentNotes: notes,
		LastMeeting: lastMeeting,
	}, nil
}

func (s *aiApp) SendFeedback(ctx context.Context, usageID int64, feedback string, comment string) error {
	return s.dm.AI().UpdateUsageFeedback(ctx, usageID, feedback, comment)
}

func (s *aiApp) GetUsageReport(ctx context.Context, period string) (entity.AIUsageReport, error) {
	userID, err := s.authApp.GetLoggedUserID(ctx)
	if err != nil {
		return entity.AIUsageReport{}, fmt.Errorf("failed to get logged user: %w", err)
	}

	return s.dm.AI().GetUsageReport(ctx, userID, period)
}

func (s *aiApp) buildContextPrompt(context entity.PersonAIContext) string {
	prompt := fmt.Sprintf("CONTEXT ABOUT %s:\n\n", context.Person.Name)

	if context.Person.Position != "" {
		prompt += fmt.Sprintf("Position: %s\n", context.Person.Position)
	}
	if context.Person.Department != "" {
		prompt += fmt.Sprintf("Department: %s\n", context.Person.Department)
	}

	if len(context.Attributes) > 0 {
		prompt += "\nPERSONAL ATTRIBUTES:\n"
		for key, value := range context.Attributes {
			prompt += fmt.Sprintf("- %s: %s\n", key, value)
		}
	}

	if len(context.RecentNotes) > 0 {
		prompt += "\nRECENT EVENTS:\n"
		count := 0
		for _, note := range context.RecentNotes {
			if count >= 5 {
				break
			}
			prompt += fmt.Sprintf("- %s: %s\n", note.CreatedAt.Format("2006-01-02"), note.Content)
			count++
		}
	}

	if context.LastMeeting != nil {
		prompt += fmt.Sprintf("\nLAST 1:1 MEETING: %s\n", context.LastMeeting.Date.Format("2006-01-02"))
	}

	return prompt
}
