package routes

import (
	"errors"
	"log"
	"net/http"
	"rminder/internal/database"
	"rminder/internal/middleware"
	"rminder/web"
	"rminder/web/components"
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

	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
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
	err = web.ListsContent(list, persistence, false).Render(ctx.Request.Context(), ctx.Writer)
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
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
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
		id, _ := strconv.Atoi(listID)
		err = db.DeleteList(id)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error deleting list. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	lists, err := db.Lists()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling DeleteList lists. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling DeleteList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e", err)
	}
}

// TODO: update this to work with lists
func UpdateList(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	err := ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %v", err)
	}

	// update new task
	id, _ := strconv.Atoi(listID)
	name := ctx.Request.FormValue("new-list")
	pinned := ctx.Request.FormValue("pin")
	swatch := ctx.Request.FormValue("swatch")
	icon := ctx.Request.FormValue("icon")
	if name != "" && len(name) >= 3 && len(name) <= 255 && swatch != "" && icon != "" {
		err := db.UpdateList(id, name, swatch, icon, pinned == "1")
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
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e", err)
	}
}

func SearchLists(ctx *gin.Context) {
	db := middleware.GetUserDatabase(ctx)

	var (
		lists []*database.List
		err   error
	)

	err = ctx.Request.ParseForm()
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	query := ctx.Request.FormValue("query")

	if query != "" && len(query) >= 3 {
		lists, err = db.SearchLists(query)

		err = db.UpdatePersistence(0, 0, 0)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating persistence list. Err: %v", err)
		}
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("Search query validation failed"))
		log.Fatalf("Search query validation failed. Err: %v", err)
	}

	if err != nil {
		log.Fatalf("error handling SearchLists. Err: %v", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling SearchLists Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.SearchListContent(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in SearchLists: %e", err)
	}
}
