package user

import (
	"github.com/gin-contrib/sessions"
)

type User struct {
	Id         string `json:"id"`
	HasPremium bool   `json:"has_premium"`
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
