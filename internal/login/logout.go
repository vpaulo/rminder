package login

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"rminder/internal/pkg/config"
)

// Handler for our logout.
func LogoutHandler(cfg config.AuthConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logoutUrl, err := url.Parse("https://" + cfg.Domain + "/v2/logout")
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session := sessions.Default(ctx)
		session.Clear()
		err = session.Save()
		if err != nil {
			log.Fatalf("Error saving session: %e", err)
		}

		parameters := url.Values{}
		parameters.Add("returnTo", cfg.ReturnUrl)
		parameters.Add("client_id", cfg.ClientID)
		logoutUrl.RawQuery = parameters.Encode()

		ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	}
}
