package authenticator

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func userIdFromProfile(profile map[string]interface{}) (string, error) {
	sub := fmt.Sprintf("%v", profile["sub"])
	hasher := sha1.New()
	_, err := hasher.Write([]byte(sub))
	if err != nil {
		return "", err
	}
	user_id := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return user_id, nil
}

// Handler for our callback.
func CallbackHandler(auth *Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		user_id, err := userIdFromProfile(profile)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		session.Set("user_id", user_id)

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)

		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
