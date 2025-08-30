package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	t := template.New("home")
	template, err := t.ParseFiles("./home.tmpl.html")
	upgrader := websocket.Upgrader{}

	if err != nil {
		log.Fatalf("template: %s", err.Error())
	}

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
			log.Println(err)
			return
		}

		go func(c *websocket.Conn) {
			for {
				err = c.WriteMessage(websocket.TextMessage, []byte("Hello!"))
				if err != nil {
					logger.Error("ws-error", "write-message", err.Error())
				}
				time.Sleep(20 * time.Second)
				logger.Info("Write hello...")
			}
		}(conn)
		time.Sleep(2500 * time.Millisecond)

	})

	srv := &http.Server{
		Addr:    "localhost:4000",
		Handler: mux,
	}
	logger.Info("server-start", "addr", srv.Addr)
	srv.ListenAndServe()
}
