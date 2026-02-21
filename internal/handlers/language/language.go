package language

import (
	"net/http"
	i18npkg "rminder/internal/i18n"

	"github.com/gin-gonic/gin"
)

func Set(ctx *gin.Context) {
	lang := ctx.Param("lang")

	supported := false
	for _, l := range i18npkg.SupportedLanguages {
		if l.Code == lang {
			supported = true
			break
		}
	}

	if !supported {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.SetCookie("lang", lang, 365*24*3600, "/", "", false, false)

	referer := ctx.GetHeader("Referer")
	if referer == "" {
		referer = "/"
	}
	ctx.Redirect(http.StatusFound, referer)
}
