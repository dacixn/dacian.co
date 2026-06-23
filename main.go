package main

import (
	"html/template"
	"net/http"
	"time"

	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

type StatsData struct {
	Memory string
	CPU    string
	Uptime string
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/stats", statsHandler)
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "projects.html", nil)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Percent(time.Second, false)
	h, _ := host.Info()

	data := StatsData{
		Memory: fmt.Sprintf("%.1f/%.1fGB", float64(v.Used)/1e9, float64(v.Total)/1e9),
		CPU:    fmt.Sprintf("%.0f%%", c[0]),
		Uptime: fmt.Sprintf("%d days", h.Uptime/86400),
	}

	templates.ExecuteTemplate(w, "stats.html", data)
}
