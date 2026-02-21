package middleware

import (
	i18npkg "rminder/internal/i18n"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang, _ := ctx.Cookie("lang")
		if lang == "" {
			lang = ctx.GetHeader("Accept-Language")
		}
		localizer := i18n.NewLocalizer(i18npkg.Bundle, lang, "en")
		ctx.Set("lang", lang)
		ctx.Set("localizer", localizer)
		ctx.Next()
	}
}
