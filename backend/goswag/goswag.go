package main

import (
	_ "github.com/diegoclair/go_utils/resterrors"
	_ "github.com/diegoclair/leaderpro/internal/transport/rest/viewmodel"
	_ "github.com/diegoclair/leaderpro/internal/transport/rest/routes/pingroute"
)

//	@Summary		Logout
//	@Description	Logout the user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user-token	header	string	true	"User access token"
//	@Success		200
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/auth/logout [post]
func handleLogout() {} //nolint:unused 

//	@Summary		Create a new company
//	@Description	Create a new company
//	@Tags			companies
//	@Accept			json
//	@Produce		json
//	@Param			request		body	viewmodel.CompanyRequest	true	"Request"
//	@Param			user-token	header	string						true	"User access token"
//	@Success		201
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies [post]
func handleCreateCompany() {} //nolint:unused 

//	@Summary		Get user companies
//	@Description	Get all companies for the authenticated user
//	@Tags			companies
//	@Produce		json
//	@Param			user-token	header		string	true	"User access token"
//	@Success		200			{object}	[]viewmodel.CompanyResponse
//	@Failure		400			{object}	resterrors.restErr
//	@Failure		404			{object}	resterrors.restErr
//	@Failure		500			{object}	resterrors.restErr
//	@Failure		401			{object}	resterrors.restErr
//	@Failure		422			{object}	resterrors.restErr
//	@Failure		409			{object}	resterrors.restErr
//	@Router			/companies [get]
func handleGetCompanies() {} //nolint:unused 

//	@Summary		Get company by UUID
//	@Description	Get company details by UUID
//	@Tags			companies
//	@Produce		json
//	@Param			company_uuid	path		string	true	"company uuid"
//	@Param			user-token		header		string	true	"User access token"
//	@Success		200				{object}	viewmodel.CompanyResponse
//	@Failure		400				{object}	resterrors.restErr
//	@Failure		404				{object}	resterrors.restErr
//	@Failure		500				{object}	resterrors.restErr
//	@Failure		401				{object}	resterrors.restErr
//	@Failure		422				{object}	resterrors.restErr
//	@Failure		409				{object}	resterrors.restErr
//	@Router			/companies/:company_uuid [get]
func handleGetCompanyByUUID() {} //nolint:unused 

//	@Summary		Update company
//	@Description	Update company by UUID
//	@Tags			companies
//	@Accept			json
//	@Produce		json
//	@Param			request			body	viewmodel.CompanyRequest	true	"Request"
//	@Param			company_uuid	path	string						true	"company uuid"
//	@Param			user-token		header	string						true	"User access token"
//	@Success		204
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies/:company_uuid [put]
func handleUpdateCompany() {} //nolint:unused 

//	@Summary		Delete company
//	@Description	Delete company by UUID
//	@Tags			companies
//	@Produce		json
//	@Param			company_uuid	path	string	true	"company uuid"
//	@Param			user-token		header	string	true	"User access token"
//	@Success		204
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies/:company_uuid [delete]
func handleDeleteCompany() {} //nolint:unused 

//	@Summary		Create a new person
//	@Description	Create a new person in the company
//	@Tags			companies/:company_uuid/people
//	@Accept			json
//	@Produce		json
//	@Param			request			body	viewmodel.PersonRequest	true	"Request"
//	@Param			company_uuid	path	string					true	"company uuid"
//	@Param			user-token		header	string					true	"User access token"
//	@Success		201
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies/:company_uuid/people [post]
func handleCreatePerson() {} //nolint:unused 

//	@Summary		Get company people
//	@Description	Get all people in the company, optionally filtered by search
//	@Tags			companies/:company_uuid/people
//	@Produce		json
//	@Param			company_uuid	path		string	true	"company uuid"
//	@Param			search			query		string	false	"search term to filter people"
//	@Param			user-token		header		string	true	"User access token"
//	@Success		200				{object}	[]viewmodel.PersonResponse
//	@Failure		400				{object}	resterrors.restErr
//	@Failure		404				{object}	resterrors.restErr
//	@Failure		500				{object}	resterrors.restErr
//	@Failure		401				{object}	resterrors.restErr
//	@Failure		422				{object}	resterrors.restErr
//	@Failure		409				{object}	resterrors.restErr
//	@Router			/companies/:company_uuid/people [get]
func handleGetCompanyPeople() {} //nolint:unused 

//	@Summary		Get person by UUID
//	@Description	Get person details by UUID
//	@Tags			companies/:company_uuid/people
//	@Produce		json
//	@Param			company_uuid	path		string	true	"company uuid"
//	@Param			person_uuid		path		string	true	"person uuid"
//	@Param			user-token		header		string	true	"User access token"
//	@Success		200				{object}	viewmodel.PersonResponse
//	@Failure		400				{object}	resterrors.restErr
//	@Failure		404				{object}	resterrors.restErr
//	@Failure		500				{object}	resterrors.restErr
//	@Failure		401				{object}	resterrors.restErr
//	@Failure		422				{object}	resterrors.restErr
//	@Failure		409				{object}	resterrors.restErr
//	@Router			/companies/:company_uuid/people/:person_uuid [get]
func handleGetPersonByUUID() {} //nolint:unused 

//	@Summary		Update person
//	@Description	Update person by UUID
//	@Tags			companies/:company_uuid/people
//	@Accept			json
//	@Produce		json
//	@Param			request			body	viewmodel.PersonRequest	true	"Request"
//	@Param			company_uuid	path	string					true	"company uuid"
//	@Param			person_uuid		path	string					true	"person uuid"
//	@Param			user-token		header	string					true	"User access token"
//	@Success		204
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies/:company_uuid/people/:person_uuid [put]
func handleUpdatePerson() {} //nolint:unused 

//	@Summary		Delete person
//	@Description	Delete person by UUID
//	@Tags			companies/:company_uuid/people
//	@Produce		json
//	@Param			company_uuid	path	string	true	"company uuid"
//	@Param			person_uuid		path	string	true	"person uuid"
//	@Param			user-token		header	string	true	"User access token"
//	@Success		204
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/companies/:company_uuid/people/:person_uuid [delete]
func handleDeletePerson() {} //nolint:unused 

//	@Summary		Get User Profile
//	@Description	Get the current user's profile
//	@Tags			users
//	@Produce		json
//	@Param			user-token	header		string	true	"User access token"
//	@Success		200			{object}	viewmodel.User
//	@Failure		400			{object}	resterrors.restErr
//	@Failure		404			{object}	resterrors.restErr
//	@Failure		500			{object}	resterrors.restErr
//	@Failure		401			{object}	resterrors.restErr
//	@Failure		422			{object}	resterrors.restErr
//	@Failure		409			{object}	resterrors.restErr
//	@Router			/users/profile [get]
func handleGetProfile() {} //nolint:unused 

//	@Summary		Update User Profile
//	@Description	Update the current user's profile
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request		body		viewmodel.UpdateUser	true	"Request"
//	@Param			user-token	header		string					true	"User access token"
//	@Success		200			{object}	viewmodel.User
//	@Failure		400			{object}	resterrors.restErr
//	@Failure		404			{object}	resterrors.restErr
//	@Failure		500			{object}	resterrors.restErr
//	@Failure		401			{object}	resterrors.restErr
//	@Failure		422			{object}	resterrors.restErr
//	@Failure		409			{object}	resterrors.restErr
//	@Router			/users/profile [put]
func handleUpdateProfile() {} //nolint:unused 

//	@Summary		Login
//	@Description	Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		viewmodel.Login	true	"Request"
//	@Success		200		{object}	viewmodel.LoginResponse
//	@Failure		400		{object}	resterrors.restErr
//	@Failure		404		{object}	resterrors.restErr
//	@Failure		500		{object}	resterrors.restErr
//	@Failure		401		{object}	resterrors.restErr
//	@Failure		422		{object}	resterrors.restErr
//	@Failure		409		{object}	resterrors.restErr
//	@Router			/auth/login [post]
func handleLogin() {} //nolint:unused 

//	@Summary		Refresh Token
//	@Description	Generate a new token using the refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		viewmodel.RefreshTokenRequest	true	"Request"
//	@Success		200		{object}	viewmodel.RefreshTokenResponse
//	@Failure		400		{object}	resterrors.restErr
//	@Failure		404		{object}	resterrors.restErr
//	@Failure		500		{object}	resterrors.restErr
//	@Failure		401		{object}	resterrors.restErr
//	@Failure		422		{object}	resterrors.restErr
//	@Failure		409		{object}	resterrors.restErr
//	@Router			/auth/refresh-token [post]
func handleRefreshToken() {} //nolint:unused 

//	@Summary		Ping the server
//	@Description	Ping the server to check if it is alive
//	@Tags			ping
//	@Produce		json
//	@Success		200	{object}	pingroute.pingResponse
//	@Failure		400	{object}	resterrors.restErr
//	@Failure		404	{object}	resterrors.restErr
//	@Failure		500	{object}	resterrors.restErr
//	@Failure		401	{object}	resterrors.restErr
//	@Failure		422	{object}	resterrors.restErr
//	@Failure		409	{object}	resterrors.restErr
//	@Router			/ping/ [get]
func handlePing() {} //nolint:unused 

//	@Summary		Create User
//	@Description	Create a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		viewmodel.CreateUser	true	"Request"
//	@Success		200		{object}	viewmodel.User
//	@Failure		400		{object}	resterrors.restErr
//	@Failure		404		{object}	resterrors.restErr
//	@Failure		500		{object}	resterrors.restErr
//	@Failure		401		{object}	resterrors.restErr
//	@Failure		422		{object}	resterrors.restErr
//	@Failure		409		{object}	resterrors.restErr
//	@Router			/users [post]
func handleCreateUser() {} //nolint:unused 

