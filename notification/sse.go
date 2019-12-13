package notification

import (
	"encoding/json"

	"github.com/ailisp/reallyfastci/core"
	"github.com/r3labs/sse"
)

var SseServer *sse.Server

func InitSse() {
	SseServer = sse.New()
	SseServer.CreateStream("build-status")
}

func NotifySse(event *core.BuildEvent) {
	data, _ := json.Marshal(map[string]string{
		"commit": event.Commit,
		"status": core.BuildStatusStr(event.Status),
	})
	SseServer.Publish("build-status", &sse.Event{
		Data: data,
	})
}
