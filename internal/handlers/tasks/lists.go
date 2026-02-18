package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"rminder/internal/app"
	"rminder/internal/app/database"
	"rminder/web"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling GetLists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling GetLists Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Render(ctx.Writer, "sidebar-lists", map[string]any{
		"Lists":       lists,
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering in GetLists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render lists.")
	}
}

func GetList(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	listID := ctx.Param("listID")

	var (
		lists []*database.List
		list  *database.List
		err   error
	)

	if listID != "" {
		list, err = db.List(listID)
		if err != nil {
			log.Error("error handling list", "listID", listID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to load list.")
			return
		}

		id, _ := strconv.Atoi(listID)
		err = db.UpdatePersistence(0, id)
		if err != nil {
			log.Error("error updating persistence list", "listID", listID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update list state.")
			return
		}

		if list.FilterBy != "" {
			lists, err = db.Lists(list.FilterBy)
			if err != nil {
				log.Error("error handling GetLists", "error", err)
				app.ErrorInternalHTML(ctx, "Failed to load filtered lists.")
				return
			}
		}
	} else {
		log.Error("error no list id")
		app.ErrorBadRequestHTML(ctx, "List ID must not be empty.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling GetList Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load list state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if list.FilterBy == "" {
		err = web.Render(ctx.Writer, "lists-content", map[string]any{
			"List":        list,
			"Persistence": persistence,
			"IsMultilist": false,
		})
	} else {
		err = web.Render(ctx.Writer, "multi-list-content", map[string]any{
			"Lists":       lists,
			"Title":       list.Name,
			"Persistence": persistence,
		})
	}
	if err != nil {
		log.Error("error rendering in ListsContent", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render list content.")
	}
}

func CreateList(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
	}

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
			log.Error("error creating list", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to create list.")
			return
		}
		log.Info("list created", "name", list)
	} else {
		log.Error("error form field validation failed", "name", list)
		app.ErrorBadRequestHTML(ctx, "All fields are required and name must be between 3 and 255 characters.")
		return
	}

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling CreateList Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load list state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "sidebar-lists", map[string]any{
		"Lists":       lists,
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering in Lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render lists.")
	}
}

func DeleteList(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	listID := ctx.Param("listID")

	var id int
	var err error

	if listID != "" {
		id, _ = strconv.Atoi(listID)
		err = db.DeleteList(id)
		if err != nil {
			log.Error("error deleting list", "listID", listID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to delete list.")
			return
		}
		log.Info("list deleted", "listID", listID)
	} else {
		log.Error("error list id is empty")
		app.ErrorBadRequestHTML(ctx, "List ID must not be empty.")
		return
	}

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling DeleteList lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling DeleteList Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load list state.")
		return
	}

	if persistence.ListId == id {
		err = db.UpdatePersistenceList(0)
		if err != nil {
			log.Error("error handling DeleteList Persistence Update", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update list state.")
			return
		}
		persistence.ListId = 0
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "sidebar-lists", map[string]any{
		"Lists":       lists,
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering in Lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render lists.")
	}
}

func UpdateList(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	listID := ctx.Param("listID")

	err := ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
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
			log.Error("error updating list", "listID", listID, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update list.")
			return
		}
		log.Info("list updated", "listID", listID, "name", name)
	} else {
		log.Error("error form field validation failed", "listID", listID)
		app.ErrorBadRequestHTML(ctx, "All fields are required and name must be between 3 and 255 characters.")
		return
	}

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error handling lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load lists.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling UpdateList Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load list state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = web.Render(ctx.Writer, "sidebar-lists", map[string]any{
		"Lists":       lists,
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering in Lists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render lists.")
	}
}

func SearchLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	var (
		lists []*database.List
		err   error
	)

	err = ctx.Request.ParseForm()
	if err != nil {
		log.Error("error parsing form", "error", err)
		app.ErrorBadRequestHTML(ctx, "Invalid form data.")
		return
	}

	query := ctx.Request.FormValue("query")

	if query != "" && len(query) >= 3 {
		lists, err = db.SearchLists(query)
		if err != nil {
			log.Error("error search list", "query", query, "error", err)
			app.ErrorInternalHTML(ctx, "Failed to search lists.")
			return
		}

		err = db.UpdatePersistence(0, 0)
		if err != nil {
			log.Error("error updating persistence list", "error", err)
			app.ErrorInternalHTML(ctx, "Failed to update search state.")
			return
		}
	} else {
		log.Error("search query validation failed", "query", query)
		app.ErrorBadRequestHTML(ctx, "Search query must be at least 3 characters.")
		return
	}

	persistence, err := db.Persistence()
	if err != nil {
		log.Error("error handling SearchLists Persistence", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to load search state.")
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = web.Render(ctx.Writer, "multi-list-content", map[string]any{
		"Lists":       lists,
		"Title":       "Search Results",
		"Persistence": persistence,
	})
	if err != nil {
		log.Error("error rendering in SearchLists", "error", err)
		app.ErrorInternalHTML(ctx, "Failed to render search results.")
	}
}

func ReorderLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)
	var reorder []database.Reorder
	var err error

	if err = ctx.ShouldBindJSON(&reorder); err != nil {
		log.Error("error binding reorder JSON", "error", err)
		app.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request data.")
		return
	}

	err = db.ReorderLists(reorder)

	if err != nil {
		log.Error("error reordering lists", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Lists order update unsuccessful.")
		return
	}

	log.Info("lists reordered")
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Lists order update successful.",
		"status":  http.StatusOK,
	})
}

func ExportLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	lists, err := db.Lists("")
	if err != nil {
		log.Error("error exporting lists", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Lists export unsuccessful.")
		return
	}

	log.Info("lists exported")
	ctx.IndentedJSON(http.StatusOK, lists)
}

func ImportLists(ctx *gin.Context) {
	db := app.GetUserDatabase(ctx)
	log := app.GetLogger(ctx)

	var err error
	var openedFile multipart.File
	var file []byte

	// Parse the multipart form, 10 MB max upload size
	ctx.Request.ParseMultipartForm(10 << 20)

	// Retrieve the file from form data
	formFile, err := ctx.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			log.Error("error no file submitted", "error", err)
			app.ErrorJSON(ctx, http.StatusBadRequest, "No file submitted.")
		} else {
			log.Error("error retrieving the file", "error", err)
			app.ErrorJSON(ctx, http.StatusBadRequest, "Error retrieving the file.")
		}
		return
	}

	openedFile, err = formFile.Open()

	if err != nil {
		log.Error("error not able to open file", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Not able to open file.")
		return
	}

	file, err = io.ReadAll(openedFile)

	if err != nil {
		log.Error("error not able to read file", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Not able to read file.")
		return
	}

	var lists []*database.List
	if err := json.Unmarshal(file, &lists); err != nil {
		log.Error("error not able to unmarshal file", "error", err)
		app.ErrorJSON(ctx, http.StatusBadRequest, "Invalid file format.")
		return
	}

	if err := db.ImportLists(lists); err != nil {
		log.Error("error failed to save imported lists", "error", err)
		app.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to save imported lists.")
		return
	}

	log.Info("lists imported", "count", len(lists))
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "Lists imported successfuly.",
		"status":  http.StatusOK,
	})
}
