package app

import (
	"rminder/internal/app/database"
	"rminder/internal/app/user"
	"rminder/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetUserDatabase(ctx *gin.Context, db database.Service) {
	ctx.Set("user_database", db)
}

func GetUserDatabase(ctx *gin.Context) database.Service {
	return ctx.MustGet("user_database").(database.Service)
}

func SetUser(ctx *gin.Context, user *user.User) {
	ctx.Set("user", user)
}

func GetUser(ctx *gin.Context) *user.User {
	return ctx.MustGet("user").(*user.User)
}

func SetLogger(ctx *gin.Context, log *logger.Logger) {
	ctx.Set("logger", log)
}

func GetLogger(ctx *gin.Context) *logger.Logger {
	return ctx.MustGet("logger").(*logger.Logger)
}

func GetCSRFToken(ctx *gin.Context) string {
	token, _ := ctx.Get("csrf_token")
	if token == nil {
		return ""
	}
	return token.(string)
}

func GetRequestID(ctx *gin.Context) string {
	return ctx.GetString("request_id")
}
