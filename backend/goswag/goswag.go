package main

import (
	_ "github.com/diegoclair/go_utils/resterrors"
	_ "github.com/diegoclair/leaderpro/internal/transport/rest/routes/pingroute"
	_ "github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
)

// @Summary		Chat with AI Leadership Coach
// @Description	Send a message to the AI leadership coach for advice and guidance
// @Tags			companies/:company_uuid/ai
// @Accept			json
// @Produce		json
// @Param			request			body		viewmodel.AIChatRequest	true	"Request"
// @Param			company_uuid	path		string					true	"Company UUID"
// @Param			user-token		header		string					true	"User access token"
// @Success		200				{object}	viewmodel.AIChatResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/ai/chat [post]
func handleChatWithAI() {} //nolint:unused

// @Summary		Get AI usage report
// @Description	Get usage statistics and costs for AI features
// @Tags			companies/:company_uuid/ai
// @Produce		json
// @Param			company_uuid	path		string	true	"Company UUID"
// @Param			period			query		string	false	"Period (today, week, month, year, all)"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.AIUsageReportResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/ai/usage [get]
func handleGetUsageReport() {} //nolint:unused

// @Summary		Send feedback about AI response
// @Description	Provide feedback (helpful/not helpful) about an AI response
// @Tags			companies/:company_uuid/ai
// @Accept			json
// @Produce		json
// @Param			request			body		viewmodel.AIFeedbackRequest	true	"Request"
// @Param			company_uuid	path		string						true	"Company UUID"
// @Param			usage_id		path		string						true	"Usage ID"
// @Param			user-token		header		string						true	"User access token"
// @Success		200				{object}	viewmodel.MessageResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/ai/usage/:usage_id/feedback [post]
func handleSendFeedback() {} //nolint:unused

// @Summary		Chat with AI about specific person
// @Description	Send a message to the AI with context about a specific person
// @Tags			companies/:company_uuid/people
// @Accept			json
// @Produce		json
// @Param			request			body		viewmodel.AIPersonChatRequest	true	"Request"
// @Param			company_uuid	path		string							true	"Company UUID"
// @Param			person_uuid		path		string							true	"Person UUID"
// @Param			user-token		header		string							true	"User access token"
// @Success		200				{object}	viewmodel.AIChatResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/ai/chat [post]
func handleChatAboutPerson() {} //nolint:unused

// @Summary		Get company by UUID
// @Description	Get company details by UUID
// @Tags			companies
// @Produce		json
// @Param			company_uuid	path		string	true	"company uuid"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.CompanyResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid [get]
func handleGetCompanyByUUID() {} //nolint:unused

// @Summary		Update company
// @Description	Update company by UUID
// @Tags			companies
// @Accept			json
// @Produce		json
// @Param			request			body	viewmodel.CompanyRequest	true	"Request"
// @Param			company_uuid	path	string						true	"company uuid"
// @Param			user-token		header	string						true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid [put]
func handleUpdateCompany() {} //nolint:unused

// @Summary		Delete company
// @Description	Delete company by UUID
// @Tags			companies
// @Produce		json
// @Param			company_uuid	path	string	true	"company uuid"
// @Param			user-token		header	string	true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid [delete]
func handleDeleteCompany() {} //nolint:unused

// @Summary		Get dashboard data
// @Description	Get dashboard data with people and statistics for a specific company
// @Tags			companies/:company_uuid/dashboard
// @Produce		json
// @Param			company_uuid	path		string	true	"Company UUID to get dashboard data"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.DashboardResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/dashboard [get]
func handleGetDashboard() {} //nolint:unused

// @Summary		Create a new person
// @Description	Create a new person in the company
// @Tags			companies/:company_uuid/people
// @Accept			json
// @Produce		json
// @Param			request			body	viewmodel.PersonRequest	true	"Request"
// @Param			company_uuid	path	string					true	"company uuid"
// @Param			user-token		header	string					true	"User access token"
// @Success		201
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people [post]
func handleCreatePerson() {} //nolint:unused

// @Summary		Get company people
// @Description	Get all people in the company, optionally filtered by search
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path		string	true	"company uuid"
// @Param			search			query		string	false	"search term to filter people"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	[]viewmodel.PersonResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people [get]
func handleGetCompanyPeople() {} //nolint:unused

// @Summary		Get person by UUID
// @Description	Get person details by UUID
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path		string	true	"company uuid"
// @Param			person_uuid		path		string	true	"person uuid"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.PersonResponse
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid [get]
func handleGetPersonByUUID() {} //nolint:unused

// @Summary		Update person
// @Description	Update person by UUID
// @Tags			companies/:company_uuid/people
// @Accept			json
// @Produce		json
// @Param			request			body	viewmodel.PersonRequest	true	"Request"
// @Param			company_uuid	path	string					true	"company uuid"
// @Param			person_uuid		path	string					true	"person uuid"
// @Param			user-token		header	string					true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid [put]
func handleUpdatePerson() {} //nolint:unused

// @Summary		Delete person
// @Description	Delete person by UUID
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path	string	true	"company uuid"
// @Param			person_uuid		path	string	true	"person uuid"
// @Param			user-token		header	string	true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid [delete]
func handleDeletePerson() {} //nolint:unused

// @Summary		Create a note for a person
// @Description	Create a new note (1:1, feedback, or observation) for a person
// @Tags			companies/:company_uuid/people
// @Accept			json
// @Produce		json
// @Param			request			body	viewmodel.CreateNoteRequest	true	"Request"
// @Param			company_uuid	path	string						true	"company uuid"
// @Param			person_uuid		path	string						true	"person uuid"
// @Param			user-token		header	string						true	"User access token"
// @Success		201
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/notes [post]
func handleCreateNote() {} //nolint:unused

// @Summary		Get person timeline
// @Description	Get timeline of direct notes for a person (1:1s and observations, excluding feedbacks/mentions)
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path		string	true	"company uuid"
// @Param			person_uuid		path		string	true	"person uuid"
// @Param			page			query		number	false	"page number"
// @Param			quantity		query		number	false	"items per page"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.PaginatedResponse[[]viewmodel.TimelineResponse]
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/timeline [get]
func handleGetPersonTimeline() {} //nolint:unused

// @Summary		Get person mentions
// @Description	Get notes where this person was mentioned (feedbacks received)
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path		string	true	"company uuid"
// @Param			person_uuid		path		string	true	"person uuid"
// @Param			page			query		number	false	"page number"
// @Param			quantity		query		number	false	"items per page"
// @Param			user-token		header		string	true	"User access token"
// @Success		200				{object}	viewmodel.PaginatedResponse[[]viewmodel.MentionResponse]
// @Failure		400				{object}	resterrors.restErr
// @Failure		404				{object}	resterrors.restErr
// @Failure		500				{object}	resterrors.restErr
// @Failure		401				{object}	resterrors.restErr
// @Failure		422				{object}	resterrors.restErr
// @Failure		409				{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/mentions [get]
func handleGetPersonMentions() {} //nolint:unused

// @Summary		Update a note
// @Description	Update an existing note (1:1, feedback, or observation) by UUID
// @Tags			companies/:company_uuid/people
// @Accept			json
// @Produce		json
// @Param			request			body	viewmodel.UpdateNoteRequest	true	"Request"
// @Param			company_uuid	path	string						true	"company uuid"
// @Param			person_uuid		path	string						true	"person uuid"
// @Param			note_uuid		path	string						true	"note uuid"
// @Param			user-token		header	string						true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/notes/:note_uuid [put]
func handleUpdateNote() {} //nolint:unused

// @Summary		Delete a note
// @Description	Delete an existing note (1:1, feedback, or observation) by UUID
// @Tags			companies/:company_uuid/people
// @Produce		json
// @Param			company_uuid	path	string	true	"company uuid"
// @Param			person_uuid		path	string	true	"person uuid"
// @Param			note_uuid		path	string	true	"note uuid"
// @Param			user-token		header	string	true	"User access token"
// @Success		204
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies/:company_uuid/people/:person_uuid/notes/:note_uuid [delete]
func handleDeleteNote() {} //nolint:unused

// @Summary		Logout
// @Description	Logout the user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			user-token	header	string	true	"User access token"
// @Success		200
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/auth/logout [post]
func handleLogout() {} //nolint:unused

// @Summary		Create a new company
// @Description	Create a new company
// @Tags			companies
// @Accept			json
// @Produce		json
// @Param			request		body	viewmodel.CompanyRequest	true	"Request"
// @Param			user-token	header	string						true	"User access token"
// @Success		201
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/companies [post]
func handleCreateCompany() {} //nolint:unused

// @Summary		Get user companies
// @Description	Get all companies for the authenticated user
// @Tags			companies
// @Produce		json
// @Param			user-token	header		string	true	"User access token"
// @Success		200			{object}	[]viewmodel.CompanyResponse
// @Failure		400			{object}	resterrors.restErr
// @Failure		404			{object}	resterrors.restErr
// @Failure		500			{object}	resterrors.restErr
// @Failure		401			{object}	resterrors.restErr
// @Failure		422			{object}	resterrors.restErr
// @Failure		409			{object}	resterrors.restErr
// @Router			/companies [get]
func handleGetCompanies() {} //nolint:unused

// @Summary		Get User Profile
// @Description	Get the current user's profile
// @Tags			users
// @Produce		json
// @Param			user-token	header		string	true	"User access token"
// @Success		200			{object}	viewmodel.User
// @Failure		400			{object}	resterrors.restErr
// @Failure		404			{object}	resterrors.restErr
// @Failure		500			{object}	resterrors.restErr
// @Failure		401			{object}	resterrors.restErr
// @Failure		422			{object}	resterrors.restErr
// @Failure		409			{object}	resterrors.restErr
// @Router			/users/profile [get]
func handleGetProfile() {} //nolint:unused

// @Summary		Update User Profile
// @Description	Update the current user's profile
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request		body		viewmodel.UpdateUser	true	"Request"
// @Param			user-token	header		string					true	"User access token"
// @Success		200			{object}	viewmodel.User
// @Failure		400			{object}	resterrors.restErr
// @Failure		404			{object}	resterrors.restErr
// @Failure		500			{object}	resterrors.restErr
// @Failure		401			{object}	resterrors.restErr
// @Failure		422			{object}	resterrors.restErr
// @Failure		409			{object}	resterrors.restErr
// @Router			/users/profile [put]
func handleUpdateProfile() {} //nolint:unused

// @Summary		Get User Preferences
// @Description	Get the current user's preferences
// @Tags			users
// @Produce		json
// @Param			user-token	header		string	true	"User access token"
// @Success		200			{object}	viewmodel.UserPreferences
// @Failure		400			{object}	resterrors.restErr
// @Failure		404			{object}	resterrors.restErr
// @Failure		500			{object}	resterrors.restErr
// @Failure		401			{object}	resterrors.restErr
// @Failure		422			{object}	resterrors.restErr
// @Failure		409			{object}	resterrors.restErr
// @Router			/users/preferences [get]
func handleGetUserPreferences() {} //nolint:unused

// @Summary		Update User Preferences
// @Description	Update the current user's preferences
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request		body		viewmodel.UpdateUserPreferences	true	"Request"
// @Param			user-token	header		string							true	"User access token"
// @Success		200			{object}	viewmodel.UserPreferences
// @Failure		400			{object}	resterrors.restErr
// @Failure		404			{object}	resterrors.restErr
// @Failure		500			{object}	resterrors.restErr
// @Failure		401			{object}	resterrors.restErr
// @Failure		422			{object}	resterrors.restErr
// @Failure		409			{object}	resterrors.restErr
// @Router			/users/preferences [put]
func handleUpdateUserPreferences() {} //nolint:unused

// @Summary		Login
// @Description	Login user and return user data with authentication tokens
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		viewmodel.Login	true	"Request"
// @Success		200		{object}	viewmodel.AuthResponse
// @Failure		400		{object}	resterrors.restErr
// @Failure		404		{object}	resterrors.restErr
// @Failure		500		{object}	resterrors.restErr
// @Failure		401		{object}	resterrors.restErr
// @Failure		422		{object}	resterrors.restErr
// @Failure		409		{object}	resterrors.restErr
// @Router			/auth/login [post]
func handleLogin() {} //nolint:unused

// @Summary		Refresh Token
// @Description	Generate a new access token using the refresh token
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		viewmodel.RefreshTokenRequest	true	"Request"
// @Success		200		{object}	viewmodel.RefreshTokenResponse
// @Failure		400		{object}	resterrors.restErr
// @Failure		404		{object}	resterrors.restErr
// @Failure		500		{object}	resterrors.restErr
// @Failure		401		{object}	resterrors.restErr
// @Failure		422		{object}	resterrors.restErr
// @Failure		409		{object}	resterrors.restErr
// @Router			/auth/refresh-token [post]
func handleRefreshToken() {} //nolint:unused

// @Summary		Ping the server
// @Description	Ping the server to check if it is alive
// @Tags			ping
// @Produce		json
// @Success		200	{object}	pingroute.pingResponse
// @Failure		400	{object}	resterrors.restErr
// @Failure		404	{object}	resterrors.restErr
// @Failure		500	{object}	resterrors.restErr
// @Failure		401	{object}	resterrors.restErr
// @Failure		422	{object}	resterrors.restErr
// @Failure		409	{object}	resterrors.restErr
// @Router			/ping/ [get]
func handlePing() {} //nolint:unused

// @Summary		Create User
// @Description	Create a new user account and return authentication tokens
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		viewmodel.CreateUser	true	"Request"
// @Success		200		{object}	viewmodel.AuthResponse
// @Failure		400		{object}	resterrors.restErr
// @Failure		404		{object}	resterrors.restErr
// @Failure		500		{object}	resterrors.restErr
// @Failure		401		{object}	resterrors.restErr
// @Failure		422		{object}	resterrors.restErr
// @Failure		409		{object}	resterrors.restErr
// @Router			/users [post]
func handleCreateUser() {} //nolint:unused
