package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

var Bundle *i18n.Bundle

type Language struct {
	Code string
	Name string
}

var SupportedLanguages = []Language{
	{Code: "en", Name: "English"},
	{Code: "pt", Name: "Português"},
}

func init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, err := localeFS.ReadDir("locales")
	if err != nil {
		panic("i18n: failed to read locales directory: " + err.Error())
	}

	for _, entry := range entries {
		data, err := localeFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			panic("i18n: failed to read locale file " + entry.Name() + ": " + err.Error())
		}
		Bundle.MustParseMessageFileBytes(data, entry.Name())
	}
}
