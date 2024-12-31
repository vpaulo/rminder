package web

import (
	"fmt"
	"strings"
	"time"
)

func taskClasses(isCompleted bool, isImportant bool, isMyDay bool) string {
	completed := ""
	important := ""
	myDay := ""

	if isCompleted {
		completed = "completed"
	}

	if isImportant {
		important = "important"
	}

	if isMyDay {
		myDay = "today"
	}

	return strings.Trim(fmt.Sprintf("%s %s %s", completed, important, myDay), " ")
}

func formatDate(date string) string {
	tm, err := time.Parse(time.RFC3339, date)

	if err != nil {
		return "NA"
	}

	return tm.Format(time.DateTime)
}
