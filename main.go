package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
)

func main() {
	http.HandleFunc("/static/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/style.css")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Message string
		}{
			Message: "hello world",
		}

		tmpl.Execute(w, data)
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().Format(time.RFC1123)))
	})

	http.HandleFunc("/mem", func(w http.ResponseWriter, r *http.Request) {
		mem, err := memory.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("%.1f/%.1fGB", float32(mem.Used)/(1024*1024*1024), float32(mem.Total)/(1024*1024*1024))))
		fmt.Println(mem.Used/(1024*1024*1024), mem.Total/(1024*1024*1024))
	})

	http.HandleFunc("/up", func(w http.ResponseWriter, r *http.Request) {
		up, err := uptime.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		days := up.Hours() / 24
		if days < 1 {
			w.Write([]byte(fmt.Sprintf("%.1f days", up.Hours()/24)))
		} else {

			w.Write([]byte(fmt.Sprintf("%.0f days", up.Hours()/24)))
		}
	})
	pages, err := Parse()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started http server on :8080")
	// http.ListenAndServe(":8080", nil)
	fmt.Println(pages["projects"].Content)
}
