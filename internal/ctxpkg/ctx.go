package ctxpkg

import (
	"ToDoWithKolya/internal/models"
	"context"
)

func UserFromContext(ctx context.Context) (models.User, bool) {
	user, ok := ctx.Value("users").(models.User)
	return user, ok
}
