package web

import (
	"embed"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	i18npkg "rminder/internal/i18n"
)

//go:embed "templates"
var templateFS embed.FS

var templates *template.Template

func init() {
	funcMap := template.FuncMap{
		"itoa":           strconv.Itoa,
		"repeatStr":      strings.Repeat,
		"formatDate":     formatDate,
		"formatDateOnly": formatDateOnly,
		"swatchColours":  swatchColours,
		"iconsList":      iconsList,
		"priorityValues": priorityValues,
		"add":            func(a, b int) int { return a + b },
		"concat": func(parts ...string) string {
			return strings.Join(parts, "")
		},
		"rawHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"dict": func(pairs ...any) map[string]any {
			m := make(map[string]any, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				m[pairs[i].(string)] = pairs[i+1]
			}
			return m
		},
	}

	templates = template.Must(
		template.New("").Funcs(funcMap).ParseFS(templateFS,
			"templates/shared/*.html",
			"templates/public/*.html",
			"templates/apps/tasks/*.html",
			"templates/apps/tasks/components/*.html",
		),
	)
}

func Render(ctx *gin.Context, name string, data map[string]any) error {
	if val, exists := ctx.Get("localizer"); exists {
		localizer := val.(*i18n.Localizer)
		data["T"] = func(key string) string {
			msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: key})
			if err != nil {
				return key
			}
			return msg
		}
	}
	data["Lang"] = ctx.GetString("lang")
	data["SupportedLanguages"] = i18npkg.SupportedLanguages
	return templates.ExecuteTemplate(ctx.Writer, name, data)
}

func formatDate(date string) string {
	tm, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}
	return tm.Format(time.DateTime)
}

func formatDateOnly(date string) string {
	tm, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}
	return tm.Format(time.DateOnly)
}

func swatchColours() []string {
	return []string{
		"--base-colour",
		"--colour-fresh-blue-500",
		"--colour-cyan-700",
		"--colour-cyan-800",
		"--colour-sunrise-yellow-1000",
		"--colour-volcano-400",
		"--colour-red-300",
		"--colour-sunset-orange-600",
		"--colour-lime-700",
		"--colour-pink-500",
		"--colour-indigo-400",
	}
}

func iconsList() []string {
	return []string{
		"list-ul-icon",
		"file-icon",
		"bars-progress-icon",
		"calendar-icon",
		"clipboard-icon",
		"clipboard-list-icon",
		"folder-icon",
		"folder-open-icon",
		"bell-icon",
		"bookmark-icon",
		"pen-icon",
	}
}

func priorityValues() []string {
	return []string{
		"None",
		"Low",
		"Medium",
		"High",
	}
}
