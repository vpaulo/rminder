package tasks

import (
	"rminder/internal/app"
	"rminder/web"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
	}

	name := ctx.Request.FormValue("name")
	if name == "" || len(name) < 3 || len(name) > 255 {
		log.Error("error form field validation failed", "name", name)
		app.ErrorBadRequestHTML(ctx, "Name must be between 3 and 255 characters.")
		return
	}

	groups, err := db.Groups()
	if err != nil {
		log.Error("error loading groups for position", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load groups.")
		return
	}

	if err := db.CreateGroup(name, len(groups)+1); err != nil {
		log.Error("error creating group", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to create group.")
		return
	}
	log.Info("group created", "name", name)

	renderSidebarLists(ctx)
}

func UpdateGroup(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	groupID := ctx.Param("groupID")
	id, _ := strconv.Atoi(groupID)

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
	}

	name := ctx.Request.FormValue("name")
	if name == "" || len(name) < 3 || len(name) > 255 {
		log.Error("error form field validation failed", "name", name)
		app.ErrorBadRequestHTML(ctx, "Name must be between 3 and 255 characters.")
		return
	}

	if err := db.UpdateGroup(id, name); err != nil {
		log.Error("error updating group", "groupID", groupID, "error", err)
		app.ErrorInternalHTML(ctx, "Failed to update group.")
		return
	}
	log.Info("group updated", "groupID", groupID, "name", name)

	renderSidebarLists(ctx)
}

func DeleteGroup(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	groupID := ctx.Param("groupID")
	id, _ := strconv.Atoi(groupID)

	if err := db.DeleteGroup(id); err != nil {
		log.Error("error deleting group", "groupID", groupID, "error", err)
		app.ErrorInternalHTML(ctx, "Failed to delete group.")
		return
	}
	log.Info("group deleted", "groupID", groupID)

	renderSidebarLists(ctx)
}

// renderSidebarLists is a shared helper to re-render the sidebar lists after group mutations.
func renderSidebarLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error loading lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	groups, err := db.Groups()
	if err != nil {
		log.Error("error loading groups", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load groups.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error loading persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "sidebar-lists", map[string]any{
		"Lists":       lists,
		"Groups":      groups,
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering sidebar-lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render lists.")
	}
}
