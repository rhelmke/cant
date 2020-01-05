package livelog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"cant/util/globals"

	"github.com/gorilla/websocket"
	//"strconv"
)

// we are using websockets :-)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 32768,
	// dirty (and in websec dangerous) hack
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Livelog implements the actual handler
func Livelog(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	log := globals.Livelog.Register(r)
	go func() {
		defer conn.Close()
		defer globals.Livelog.Unregister(r)
		// heap allocation optimization to remove pressure from the GC
		var buf bytes.Buffer
		reuse := globals.Livelog.NewLogData()
		// -----------------------------------------------------------
		for {
			select {
			case data := <-log:
				buf.Reset()
				enc := json.NewEncoder(&buf)
				reuse.Identifier = data.Identifier
				reuse.Timestamp = data.Timestamp
				reuse.Msg = data.Msg
				if err := enc.Encode(&reuse); err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
					return
				}
			}
		}
	}()
	go func() {
		defer conn.Close()
		defer globals.Livelog.Unregister(r)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}
