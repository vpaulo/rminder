package login

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"rminder/internal/app"
	"rminder/internal/app/user"
	"rminder/internal/login/authenticator"
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
func CallbackHandler(s *app.App, auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := s.Logger()

		if err := auth.Init(); err != nil {
			log.Error("failed to initialize authenticator", "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to process authentication.")
			return
		}

		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			log.Error("invalid state parameter in callback")
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			log.Error("failed to exchange authorization code", "error", err)
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			log.Error("failed to verify ID token", "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			log.Error("failed to parse ID token claims", "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to process authentication.")
			return
		}

		user_id, err := userIdFromProfile(profile)
		if err != nil {
			log.Error("failed to generate user ID from profile", "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to process authentication.")
			return
		}
		user.SetUserId(session, user_id)

		user_obj, err := s.GetUser(user_id)
		if err != nil {
			user_obj = &user.User{
				Id:         user_id,
				Name:       fmt.Sprintf("%v", profile["name"]),
				GivenName:  fmt.Sprintf("%v", profile["given_name"]),
				FamilyName: fmt.Sprintf("%v", profile["family_name"]),
				PictureUrl: fmt.Sprintf("%v", profile["picture"]),
				HasPremium: false,
			}
			log.Info("new user registered", "user_id", user_id)
		}

		err = s.SaveUser(user_obj)
		if err != nil {
			log.Error("failed to save user", "user_id", user_id, "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to save user.")
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)

		if err := session.Save(); err != nil {
			log.Error("failed to save session", "user_id", user_id, "error", err)
			ctx.String(http.StatusInternalServerError, "Failed to save session.")
			return
		}

		log.Info("user logged in", "user_id", user_id)
		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/tasks")
	}
}
