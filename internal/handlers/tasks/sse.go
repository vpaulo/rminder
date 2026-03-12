package tasks

import (
	"encoding/json"
	"fmt"
	"rminder/internal/app"
	"time"

	"github.com/gin-gonic/gin"
)

func SSEHandler(ctx *gin.Context) {
	broker := app.GetSSEBroker(ctx)
	user := app.GetUser(ctx)

	ch := broker.Subscribe(user.Id)
	defer broker.Unsubscribe(user.Id, ch)

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Writer.WriteHeader(200)
	ctx.Writer.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			fmt.Fprintf(ctx.Writer, "data: %s\n\n", data)
			ctx.Writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(ctx.Writer, ": keepalive\n\n")
			ctx.Writer.Flush()
		case <-ctx.Request.Context().Done():
			return
		}
	}
}
