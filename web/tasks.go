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

func swatchColours() []string {
	colours := []string{
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

	return colours
}

func iconsList() []string {
	icons := []string{
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

	return icons
}

func priorityValues() []string {
	// None=0, Low=1, Medium=2, High=3
	values := []string{
		"None",
		"Low",
		"Medium",
		"High",
	}

	return values
}
