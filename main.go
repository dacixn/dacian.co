package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	PORT = ":8080"
)

func main() {
	godotenv.Load()
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.LoadHTMLGlob("templates/**/*")
	r.Static("static/", "./static/")
	// r.Static("data/", "./data/")

	pages, err := getCardPagesFromJson("./data/cards.json")
	if err != nil {
		log.Println(err)
	}
	projects, err := getCardPageByName(pages, "projects")
	if err != nil {
		log.Println(err)
	}
	software, err := getCardPageByName(pages, "software")
	if err != nil {
		log.Println(err)
	}
	ct := &CurrentTrack{}
	cs := &CurrentStats{}

	fetchTrack := func() {
		tracks, err := getScrobbles(os.Getenv("LASTFM_USERNAME"), os.Getenv("LASTFM_API_KEY"))
		if err != nil {
			log.Println("failed to fetch scrobbles:", err)
			return
		}
		if len(tracks) == 0 {
			return
		}
		track := GetTrackInfo(tracks[0])
		ct.Set(track)
	}
	go func() {
		fetchTrack()
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			fetchTrack()
		}
	}()
	fetchStats := func() {
		mem, err := getMemory()
		if err != nil {
			log.Println(err)
		}
		cpu, err := getCPU()
		if err != nil {
			log.Println(err)
		}
		up, err := getUptime()
		if err != nil {
			log.Println(err)
		}
		stats := ServerStats{
			Memory: mem,
			CPU:    cpu,
			Uptime: up,
		}
		cs.Set(stats)
	}
	go func() {
		fetchStats()
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			fetchStats()
		}
	}()

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.GET("/software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cards.html", gin.H{
			"Name":        software.Name,
			"Description": software.Description,
			"Cards":       software.Cards,
		})
	})
	r.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cards.html", gin.H{
			"Name":        projects.Name,
			"Description": projects.Description,
			"Cards":       projects.Cards,
		})
	})

	r.GET("/music", func(c *gin.Context) {
		track := ct.Get()

		c.HTML(http.StatusOK, "music.html", gin.H{
			"Title":    strings.ToLower(truncate(track.Title, 25)),
			"Artist":   strings.ToLower(truncate(track.Artist, 25)),
			"ImageURL": track.ImageURL,
		})
	})
	r.GET("/stats", func(c *gin.Context) {
		stats := cs.Get()

		c.HTML(http.StatusOK, "stats.html", gin.H{
			"Memory": stats.Memory,
			"CPU":    stats.CPU,
			"Uptime": stats.Uptime,
		})
	})
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
