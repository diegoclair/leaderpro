package logger

import (
	"context"

	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/leaderpro/infra"
)

func NewLogger(appName string, debugLevel bool) logger.Logger {
	params := logger.LogParams{
		AppName:                  appName,
		DebugLevel:               debugLevel,
		AddAttributesFromContext: addDefaultAttributesToLogger,
	}
	return logger.New(params)
}

func addDefaultAttributesToLogger(ctx context.Context) []logger.LogField {
	args := []logger.LogField{}

	if sessionCode, ok := getSession(ctx); ok {
		args = append(args, logger.String("session", sessionCode))
	}

	if accountUUID, ok := getAccountUUID(ctx); ok {
		args = append(args, logger.String("account_uuid", accountUUID))
	}

	return args
}

func getContextValue(ctx context.Context, key infra.Key) string {
	if ctx == nil {
		return ""
	}

	value := ctx.Value(key)
	if value == nil {
		return ""
	}

	return value.(string)
}

func getSession(ctx context.Context) (string, bool) {
	sessionCode := getContextValue(ctx, infra.SessionKey)
	if sessionCode == "" {
		return "", false
	}

	return sessionCode, true
}

func getAccountUUID(ctx context.Context) (string, bool) {
	accountUUID := getContextValue(ctx, infra.AccountUUIDKey)
	if accountUUID == "" {
		return "", false
	}

	return accountUUID, true
}
