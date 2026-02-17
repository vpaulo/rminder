package user

import (
	"github.com/gin-contrib/sessions"
)

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	PictureUrl string `json:"picture_url"`
}

func GetUserId(session sessions.Session) string {
	user_id := session.Get("user_id")
	if user_id == nil {
		return ""
	}
	user_id_str := user_id.(string)
	return user_id_str
}

func SetUserId(session sessions.Session, user_id string) {
	session.Set("user_id", user_id)
}
