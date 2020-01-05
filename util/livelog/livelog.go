package livelog

import (
	"net/http"
)

type LogData struct {
	Identifier string `json:"identifier"`
	Timestamp  int64  `json:"timestamp"`
	Msg        string `json:"msg"`
}

type Livelog struct {
	clients    map[*http.Request]chan LogData
	Send       chan LogData
	add        chan *http.Request
	remove     chan *http.Request
	confirmAdd chan bool
}

func New() *Livelog {
	log := &Livelog{
		clients:    make(map[*http.Request]chan LogData),
		Send:       make(chan LogData, 1000),
		remove:     make(chan *http.Request, 1000),
		add:        make(chan *http.Request, 1000),
		confirmAdd: make(chan bool, 1000),
	}
	go log.Run()
	return log
}

func (log *Livelog) NewLogData() LogData {
	return LogData{}
}

func (log *Livelog) Register(client *http.Request) chan LogData {
	if _, ok := log.clients[client]; ok {
		return nil
	}
	log.add <- client
	<-log.confirmAdd
	return log.clients[client]
}

func (log *Livelog) Unregister(client *http.Request) {
	if _, ok := log.clients[client]; !ok {
		return
	}
	log.remove <- client
}

func (log *Livelog) Run() {
	for {
		select {
		case client := <-log.add:
			log.clients[client] = make(chan LogData, 1000)
			log.confirmAdd <- true
		case client := <-log.remove:
			close(log.clients[client])
			delete(log.clients, client)
		case message := <-log.Send:
			for _, client := range log.clients {
				client <- message
			}
		}
	}
}
