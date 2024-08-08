package web

import (
	"fmt"
	"strings"
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
