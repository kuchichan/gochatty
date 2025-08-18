package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	t := template.New("home")
	logger := log.New(os.Stdout, "server:", 1)
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
					logger.Fatalf("write: %s", err.Error())
				}
				time.Sleep(500 * time.Millisecond)
				logger.Println("Write hello...")
			}
		}(conn)
		time.Sleep(2500 * time.Millisecond)

	})

	srv := &http.Server{
		Addr:    "localhost:4000",
		Handler: mux,
	}

	srv.ListenAndServe()
}
