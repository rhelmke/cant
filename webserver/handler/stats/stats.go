// Package stats implements a statistics handler
package stats

import (
	"cant/util/globals"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// we are using websockets :-)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// dirty (and in websec dangerous) hack
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Stats implements the actual handler
func Stats(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer conn.Close()
		defer ticker.Stop()
		snap := globals.Statistics.Snapshot()
		if err := conn.WriteMessage(websocket.TextMessage, snap); err != nil {
			return
		}
		for {
			select {
			case <-ticker.C:
				snap := globals.Statistics.Snapshot()
				if err := conn.WriteMessage(websocket.TextMessage, snap); err != nil {
					return
				}
			}
		}
	}()
	go func() {
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}
