package login

import (
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"rminder/internal/pkg/config"
	"rminder/internal/pkg/logger"
)

// Handler for our logout.
func LogoutHandler(cfg config.AuthConfig, log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logoutUrl, err := url.Parse("https://" + cfg.Domain + "/v2/logout")
		if err != nil {
			log.Error("error parsing logout URL", "error", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session := sessions.Default(ctx)
		session.Clear()
		err = session.Save()
		if err != nil {
			log.Error("error saving session on logout", "error", err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		log.Info("user logged out")

		parameters := url.Values{}
		parameters.Add("returnTo", cfg.ReturnUrl)
		parameters.Add("client_id", cfg.ClientID)
		logoutUrl.RawQuery = parameters.Encode()

		ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	}
}
