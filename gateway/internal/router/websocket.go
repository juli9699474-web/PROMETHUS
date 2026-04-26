package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func swarmWebSocketHandler(deps *Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		sub := deps.Hub.Subscribe()
		defer deps.Hub.Unsubscribe(sub)

		for evt := range sub {
			if err := conn.WriteJSON(evt); err != nil {
				return
			}
		}
	}
}
