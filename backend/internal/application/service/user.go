package service

import (
	"context"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
	"github.com/diegoclair/leaderpro/infra"
	"github.com/diegoclair/leaderpro/internal/domain"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
	"github.com/twinj/uuid"
)

type userService struct {
	cache     contract.CacheManager
	crypto    contract.Crypto
	dm        contract.DataManager
	log       logger.Logger
	validator validator.Validator
}

func newUserSvc(infra domain.Infrastructure) contract.UserApp {
	return &userService{
		cache:     infra.CacheManager(),
		crypto:    infra.Crypto(),
		dm:        infra.DataManager(),
		log:       infra.Logger(),
		validator: infra.Validator(),
	}
}

func (s *userService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Generate UUID for the user
	user.UUID = uuid.NewV4().String()

	// Hash the password
	hashedPassword, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		s.log.Errorw(ctx, "error hashing password", logger.Err(err))
		return user, resterrors.NewInternalServerError("error processing user data")
	}
	user.Password = hashedPassword

	// Set default values
	if user.Plan == "" {
		user.Plan = "trial"
	}
	if !user.EmailVerified {
		user.EmailVerified = false
	}

	// Check if email already exists
	_, err = s.dm.User().GetUserByEmail(ctx, user.Email)
	if err == nil {
		s.log.Errorw(ctx, "email already exists", logger.String("email", user.Email))
		return user, resterrors.NewBadRequestError("email already exists")
	}
	if !mysqlutils.SQLNotFound(err.Error()) {
		s.log.Errorw(ctx, "error checking if email exists", logger.Err(err))
		return user, err
	}

	// Create user in database
	userID, err := s.dm.User().CreateUser(ctx, user)
	if err != nil {
		s.log.Errorw(ctx, "error creating user", logger.Err(err))
		return user, err
	}

	user.ID = userID

	s.log.Infow(ctx, "user created successfully",
		logger.Int64("user_id", user.ID),
		logger.String("email", user.Email),
		logger.String("name", user.Name),
	)

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	user, err := s.dm.User().GetUserByEmail(ctx, email)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return user, resterrors.NewNotFoundError("user not found")
		}
		s.log.Errorw(ctx, "error getting user by email", logger.Err(err))
		return user, err
	}

	return user, nil
}

func (s *userService) GetUserByUUID(ctx context.Context, userUUID string) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	user, err := s.dm.User().GetUserByUUID(ctx, userUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return user, resterrors.NewNotFoundError("user not found")
		}
		s.log.Errorw(ctx, "error getting user by UUID", logger.Err(err))
		return user, err
	}

	// Remove password from response
	user.Password = ""

	return user, nil
}

func (s *userService) GetLoggedUser(ctx context.Context) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userUUID := ctx.Value(infra.UserUUIDKey)
	if userUUID == nil {
		s.log.Error(ctx, "user UUID not found in context")
		return entity.User{}, resterrors.NewUnauthorizedError("user not authenticated")
	}

	user, err := s.GetUserByUUID(ctx, userUUID.(string))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) GetLoggedUserID(ctx context.Context) (int64, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userUUID := ctx.Value(infra.UserUUIDKey)
	if userUUID == nil {
		s.log.Error(ctx, "user UUID not found in context")
		return 0, resterrors.NewUnauthorizedError("user not authenticated")
	}

	userID, err := s.dm.User().GetUserIDByUUID(ctx, userUUID.(string))
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return 0, resterrors.NewNotFoundError("user not found")
		}
		s.log.Errorw(ctx, "error getting user ID by UUID", logger.Err(err))
		return 0, err
	}

	return userID, nil
}

func (s *userService) UpdateUser(ctx context.Context, userUUID string, user entity.User) error {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	// Get user ID from UUID
	userID, err := s.dm.User().GetUserIDByUUID(ctx, userUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return resterrors.NewNotFoundError("user not found")
		}
		s.log.Errorw(ctx, "error getting user ID", logger.Err(err))
		return err
	}

	// Update user
	err = s.dm.User().UpdateUser(ctx, userID, user)
	if err != nil {
		s.log.Errorw(ctx, "error updating user", logger.Err(err))
		return err
	}

	s.log.Infow(ctx, "user updated successfully",
		logger.Int64("user_id", userID),
		logger.String("user_uuid", userUUID),
	)

	return nil
}

func (s *userService) GetProfile(ctx context.Context) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	user, err := s.GetLoggedUser(ctx)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, user entity.User) (entity.User, error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	userUUID := ctx.Value(infra.UserUUIDKey)
	if userUUID == nil {
		s.log.Error(ctx, "user UUID not found in context")
		return user, resterrors.NewUnauthorizedError("user not authenticated")
	}

	// Update user profile
	err := s.UpdateUser(ctx, userUUID.(string), user)
	if err != nil {
		return user, err
	}

	// Return updated user
	updatedUser, err := s.GetUserByUUID(ctx, userUUID.(string))
	if err != nil {
		return user, err
	}

	return updatedUser, nil
}
