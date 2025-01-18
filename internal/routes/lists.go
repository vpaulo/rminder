package routes

import (
	"errors"
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/internal/middleware"
	"rminder/web"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLists(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	lists, err := db.Lists()
	if err != nil {
		log.Fatalf("error handling GetLists. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling GetLists Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.ListsContainer(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in GetLists: %e", err)
	}
}

func GetList(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	var (
		list *database.List
		err  error
	)

	if listID != "" {
		list, err = db.List(listID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error handling task. Err: %v", err)
		}

		id, _ := strconv.Atoi(listID)
		err = db.UpdatePersistence(0, id, 0)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating persistence list. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling GetList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.ListsContent(list, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in ListsContent: %e", err)
	}
}

func CreateList(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %v", err)
	}

	// create new task
	list := ctx.Request.FormValue("new-list")
	pos, e := strconv.Atoi(ctx.Request.FormValue("position"))
	pinned := ctx.Request.FormValue("pin")

	if e != nil {
		pos = 0
	}

	swatch := ctx.Request.FormValue("swatch")
	icon := ctx.Request.FormValue("icon")
	if list != "" && len(list) >= 3 && len(list) <= 255 && pos != 0 && swatch != "" && icon != "" {
		err := db.CreateList(list, swatch, icon, pos, pinned == "1")
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error creating list. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error form field validation failed. Err: %v", err)
	}

	lists, err := db.Lists()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling lists. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling CreateList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.ListsContainer(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e", err)
	}
}

func DeleteList(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	var err error

	if listID != "" {
		err = db.DeleteTask(listID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error deleting task. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = ctx.Writer.Write([]byte(""))
}

// TODO: update this to work with lists
func UpdateList(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	listID := ctx.Param("listID")
	slug := ctx.Request.PathValue("slug")

	err := ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	if listID != "" {
		var err error

		switch slug {
		case "title":
			title := ctx.Request.FormValue("title")
			if title != "" && len(title) >= 3 && len(title) <= 255 {
				err = db.UpdateTask(listID, title)
			} else {
				// TODO use ApiError here
				ctx.AbortWithError(http.StatusInternalServerError, errors.New("title validation failed"))
				log.Fatalf("error title validation failed. Err: %v", err)
			}
		case "description":
			err = db.UpdateTaskDescription(listID, ctx.Request.FormValue("description"))
		case "important":
			err = db.ToggleImportant(listID)
		case "completed":
			err = db.ToggleComplete(listID)
		case "my-day":
			// err = db.ToggleMyDay(listID)
		}

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating task. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling UpdateList Persistence. Err: %v", err)
	}

	if slug == "description" {
		// TODO return proper message
		_, _ = ctx.Writer.Write([]byte("Updated description"))
	} else {
		// get task
		task, err := db.Task(listID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error handling task. Err: %v", err)
		}

		ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = web.Task(task, persistence.TaskId).Render(ctx.Request.Context(), ctx.Writer)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("Error rendering in Task: %e", err)
		}
	}
}
