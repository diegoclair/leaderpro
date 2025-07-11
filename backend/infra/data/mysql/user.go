package mysql

import (
	"context"

	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/leaderpro/internal/domain/contract"
	"github.com/diegoclair/leaderpro/internal/domain/entity"
)

type userRepo struct {
	db dbConn
}

func newUserRepo(db dbConn) contract.UserRepo {
	return &userRepo{
		db: db,
	}
}

const userSelectBase string = `
	SELECT 
		u.user_id,
		u.user_uuid,
		u.email,
		u.name,
		u.password,
		u.phone,
		u.profile_photo,
		u.plan,
		u.trial_ends_at,
		u.subscribed_at,
		u.created_at,
		u.updated_at,
		u.last_login_at,
		u.active,
		u.email_verified
	
	FROM tab_user u
`

func (r *userRepo) parseUser(row scanner) (user entity.User, err error) {
	err = row.Scan(
		&user.ID,
		&user.UUID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Phone,
		&user.ProfilePhoto,
		&user.Plan,
		&user.TrialEndsAt,
		&user.SubscribedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.Active,
		&user.EmailVerified,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepo) CreateUser(ctx context.Context, user entity.User) (createdID int64, err error) {
	query := `
		INSERT INTO tab_user (
			user_uuid,
			email,
			name,
			password,
			phone,
			profile_photo,
			plan,
			trial_ends_at,
			subscribed_at,
			active,
			email_verified
		) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		user.UUID,
		user.Email,
		user.Name,
		user.Password,
		user.Phone,
		user.ProfilePhoto,
		user.Plan,
		user.TrialEndsAt,
		user.SubscribedAt,
		user.Active,
		user.EmailVerified,
	)
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	createdID, err = result.LastInsertId()
	if err != nil {
		return createdID, mysqlutils.HandleMySQLError(err)
	}

	return createdID, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (user entity.User, err error) {
	query := userSelectBase + `
		WHERE u.email = ? AND u.active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return user, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)
	user, err = r.parseUser(row)
	if err != nil {
		return user, mysqlutils.HandleMySQLError(err)
	}

	return user, nil
}

func (r *userRepo) GetUserByUUID(ctx context.Context, userUUID string) (user entity.User, err error) {
	query := userSelectBase + `
		WHERE u.user_uuid = ? AND u.active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return user, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userUUID)
	user, err = r.parseUser(row)
	if err != nil {
		return user, mysqlutils.HandleMySQLError(err)
	}

	return user, nil
}

func (r *userRepo) GetUserIDByUUID(ctx context.Context, userUUID string) (userID int64, err error) {
	query := `
		SELECT user_id
		FROM tab_user
		WHERE user_uuid = ? AND active = 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return userID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, userUUID).Scan(&userID)
	if err != nil {
		return userID, mysqlutils.HandleMySQLError(err)
	}

	return userID, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, userID int64, user entity.User) (err error) {
	query := `
		UPDATE tab_user
		SET 
			email = ?,
			name = ?,
			phone = ?,
			profile_photo = ?,
			plan = ?,
			trial_ends_at = ?,
			subscribed_at = ?,
			email_verified = ?,
			updated_at = NOW()
		WHERE user_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		user.Email,
		user.Name,
		user.Phone,
		user.ProfilePhoto,
		user.Plan,
		user.TrialEndsAt,
		user.SubscribedAt,
		user.EmailVerified,
		userID,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}

func (r *userRepo) UpdateLastLogin(ctx context.Context, userID int64) (err error) {
	query := `
		UPDATE tab_user
		SET 
			last_login_at = NOW(),
			updated_at = NOW()
		WHERE user_id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
}
