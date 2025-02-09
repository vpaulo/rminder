package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rminder/internal/app/database"
	"rminder/web"
	"rminder/web/components"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetLists(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	lists, err := db.Lists("")
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
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in GetLists: %e :: %v", err, e)
	}
}

func GetList(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	var (
		lists []*database.List
		list  *database.List
		err   error
	)

	if listID != "" {
		list, err = db.List(listID)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error handling task. Err: %v :: %v", err, e)
		}

		id, _ := strconv.Atoi(listID)
		err = db.UpdatePersistence(0, id, 0)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating persistence list. Err: %v :: %v", err, e)
		}

		if list.FilterBy != "" {
			lists, err = db.Lists(list.FilterBy)
			if err != nil {
				log.Fatalf("error handling GetLists. Err: %v", err)
			}
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error no list id. Err: %v :: %v", err, e)
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling GetList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if list.FilterBy == "" {
		err = web.ListsContent(list, persistence, false).Render(ctx.Request.Context(), ctx.Writer)
	} else {
		err = web.MultiListContent(lists, list.Name, persistence).Render(ctx.Request.Context(), ctx.Writer)
	}
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in ListsContent: %e :: %v", err, e)
	}
}

func CreateList(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %v :: %v", err, e)
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

	// Filters
	filter := ""
	include := ctx.Request.FormValue("include")
	completed := ctx.Request.FormValue("completed")
	important := ctx.Request.FormValue("important")
	priority := ctx.Request.FormValue("priority")
	date := ctx.Request.FormValue("date")

	// TODO: error handling
	from, _ := time.Parse("2006-01-02", ctx.Request.FormValue("from"))
	to, _ := time.Parse("2006-01-02", ctx.Request.FormValue("to"))

	if include != "" {
		filter = fmt.Sprintf("include=%s;completed=%s;important=%s;priority=%s;date=%s;from=%s;to=%s", include, completed, important, priority, date, from.Format(time.DateTime), to.Format(time.DateTime))
	}

	if list != "" && len(list) >= 3 && len(list) <= 255 && pos != 0 && swatch != "" && icon != "" {
		err := db.CreateList(list, swatch, icon, pos, pinned == "1", filter)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error creating list. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error form field validation failed. Err: %e :: %v", err, e)
	}

	lists, err := db.Lists("")
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling lists. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling CreateList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e :: %v", err, e)
	}
}

func DeleteList(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	var err error

	if listID != "" {
		id, _ := strconv.Atoi(listID)
		err = db.DeleteList(id)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error deleting list. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error list id. Err: %e :: %v", err, e)
	}

	lists, err := db.Lists("")
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling DeleteList lists. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling DeleteList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e :: %v", err, e)
	}
}

func UpdateList(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	listID := ctx.Param("listID")

	err := ctx.Request.ParseForm()
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %e :: %v", err, e)
	}

	// update new task
	id, _ := strconv.Atoi(listID)
	name := ctx.Request.FormValue("new-list")
	pinned := ctx.Request.FormValue("pin")
	swatch := ctx.Request.FormValue("swatch")
	icon := ctx.Request.FormValue("icon")

	// Filters
	filter := ""
	include := ctx.Request.FormValue("include")
	completed := ctx.Request.FormValue("completed")
	important := ctx.Request.FormValue("important")
	priority := ctx.Request.FormValue("priority")
	date := ctx.Request.FormValue("date")

	// TODO: error handling
	from, _ := time.Parse("2006-01-02", ctx.Request.FormValue("from"))
	to, _ := time.Parse("2006-01-02", ctx.Request.FormValue("to"))

	if include != "" {
		filter = fmt.Sprintf("include=%s;completed=%s;important=%s;priority=%s;date=%s;from=%s;to=%s", include, completed, important, priority, date, from.Format(time.DateTime), to.Format(time.DateTime))
	}

	if name != "" && len(name) >= 3 && len(name) <= 255 && swatch != "" && icon != "" {
		err := db.UpdateList(id, name, swatch, icon, pinned == "1", filter)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error creating list. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error form field validation failed. Err: %e :: %v", err, e)
	}

	lists, err := db.Lists("")
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("error handling lists. Err: %e :: %v", err, e)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling CreateList Persistence. Err: %v", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = components.SidebarLists(lists, persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusInternalServerError, err)
		log.Fatalf("Error rendering in Lists: %e :: %v", err, e)
	}
}

func SearchLists(ctx *gin.Context) {
	db := GetUserDatabase(ctx)

	var (
		lists []*database.List
		err   error
	)

	err = ctx.Request.ParseForm()
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("error parsing form. Err: %e :: %v", err, e)
	}

	query := ctx.Request.FormValue("query")

	if query != "" && len(query) >= 3 {
		lists, err = db.SearchLists(query)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error search list. Err: %e :: %v", err, e)
		}

		err = db.UpdatePersistence(0, 0, 0)
		if err != nil {
			e := ctx.AbortWithError(http.StatusInternalServerError, err)
			log.Fatalf("error updating persistence list. Err: %e :: %v", err, e)
		}
	} else {
		e := ctx.AbortWithError(http.StatusInternalServerError, errors.New("Search query validation failed"))
		log.Fatalf("Search query validation failed. Err: %e :: %v", err, e)
	}

	if err != nil {
		log.Fatalf("error handling SearchLists. Err: %e", err)
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Fatalf("error handling SearchLists Persistence. Err: %e", err)
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.MultiListContent(lists, "Search Results", persistence).Render(ctx.Request.Context(), ctx.Writer)
	if err != nil {
		e := ctx.AbortWithError(http.StatusBadRequest, err)
		log.Fatalf("Error rendering in SearchLists: %e :: %v", err, e)
	}
}
