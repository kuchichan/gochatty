package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn   *websocket.Conn
	toSend chan []byte
	logger *slog.Logger
	room   *Room
}

type Room struct {
	broadcast chan PayloadDecoder
	logger    *slog.Logger
}

func NewRoom(logger *slog.Logger) *Room {
	return &Room{
		broadcast: make(chan PayloadDecoder),
		logger:    logger,
	}
}

func (r *Room) run() {
	for {
		select {
		case msg := <-r.broadcast:
			r.logger.Info("room", "broadcast", msg)
		}
	}

}

type Message struct {
	MsgType string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type PayloadDecoder interface {
	DecodeMsg()
}

type JoinMsg struct {
	Name string `json:"userName"`
}

func (m *Message) DecodeMessage(data []byte) (PayloadDecoder, error) {
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	var payload PayloadDecoder

	switch m.MsgType {
	case "join":
		payload = new(JoinMsg)
		json.Unmarshal(m.Payload, payload)
	}

	return payload, nil
}

func (jm *JoinMsg) DecodeMsg() {
}

func (c *WsClient) readLoop() {
	for {
		_, message, err := c.conn.ReadMessage()
		message = bytes.TrimSpace(message)
		msg := new(Message)

		c.logger.Info("ws", "received", string(message))
		payload, err := msg.DecodeMessage(message)
		if err != nil {
			c.logger.Error("decode", "decoder-error", err.Error())
		}
		c.logger.Info("ws", "received-payload", payload)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("ws-error", "read-message", err)
			}
			break
		}

		c.room.broadcast <- payload
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t := template.New("home")
	template, err := t.ParseFiles("./home.tmpl.html")
	if err != nil {
		log.Fatalf("template: %s", err.Error())
	}

	upgrader := websocket.Upgrader{}

	room := NewRoom(logger)
	go room.run()

	fileServer := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, nil)
		if err != nil {
			log.Fatalf("execute: %s", err.Error())
		}
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			logger.Error("ws-error", "connection-error", err)
			return
		}
		client := WsClient{conn: conn, logger: logger, toSend: make(chan []byte, 256), room: room}

		go client.readLoop()

	})

	srv := &http.Server{
		Addr:    "localhost:4000",
		Handler: mux,
	}
	logger.Info("server-start", "addr", srv.Addr)
	srv.ListenAndServe()
}
