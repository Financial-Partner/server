package contextutil

import (
	"context"
)

type ContextKey string

const (
	UserEmailKey ContextKey = "user_email"
)

func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}
