package service

import (
	"context"
	"time"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/internal/application/dto"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

const (
	wrongLogin         string = "Document or password are wrong"
	errDeactivatedUser string = "User is deactivated"
)

type authApp struct {
	cache               contract.CacheManager
	crypto              contract.Crypto
	dm                  contract.DataManager
	log                 logger.Logger
	validator           validator.Validator
	userSvc             contract.UserApp
	accessTokenDuration time.Duration
}

func newAuthApp(infra domain.Infrastructure, userSvc contract.UserApp, accessTokenDuration time.Duration) *authApp {
	return &authApp{
		cache:               infra.CacheManager(),
		crypto:              infra.Crypto(),
		dm:                  infra.DataManager(),
		log:                 infra.Logger(),
		validator:           infra.Validator(),
		userSvc:             userSvc,
		accessTokenDuration: accessTokenDuration,
	}
}

func (s *authApp) Login(ctx context.Context, input dto.LoginInput) (user entity.User, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return user, err
	}

	user, err = s.dm.User().GetUserByEmail(ctx, input.Email)
	if err != nil {
		s.log.Errorw(ctx, "error getting account by document", logger.Err(err))
		return user, resterrors.NewUnauthorizedError(wrongLogin)
	}

	ctx = context.WithValue(ctx, infra.UserUUIDKey, user.UUID) // set user uuid in context to be used in logs

	if !user.Active {
		s.log.Error(ctx, "user is not active")
		return user, resterrors.NewUnauthorizedError(errDeactivatedUser)
	}

	s.log.Infow(ctx, "user information used to login",
		logger.Int64("user_id", user.ID),
		logger.String("name", user.Name),
	)

	err = s.crypto.CheckPassword(input.Password, user.Password)
	if err != nil {
		s.log.Error(ctx, "wrong password")
		return user, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return user, nil
}

func (s *authApp) CreateSession(ctx context.Context, session dto.Session) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	err = session.Validate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	_, err = s.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		s.log.Errorw(ctx, "error creating session", logger.Err(err))
		return err
	}

	return nil
}

func (s *authApp) GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	session, err = s.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return session, resterrors.NewUnauthorizedError("session not found")
		}
		s.log.Errorw(ctx, "error getting session", logger.Err(err))
		return session, err
	}

	return session, nil
}

func (s *authApp) Logout(ctx context.Context, accessToken string) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	loggedUserID, err := s.userSvc.GetLoggedUserID(ctx)
	if err != nil {
		return err
	}

	// access token will be on cache for 3 minutes after it duration
	// this is to avoid the user to login again with the same access token (used in the middleware)
	err = s.cache.SetStringWithExpiration(ctx, accessToken, "true", s.accessTokenDuration+3*time.Minute)
	if err != nil {
		s.log.Errorw(ctx, "error logging out", logger.Err(err))
		return err
	}

	err = s.dm.Auth().SetSessionAsBlocked(ctx, loggedUserID)
	if err != nil {
		s.log.Errorw(ctx, "error logging out", logger.Err(err))
		return err
	}

	return nil
}
