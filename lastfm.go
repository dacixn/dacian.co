package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Image struct {
	Size string `json:"size"`
	URL  string `json:"#text"`
}
type Artist struct {
	Name string `json:"#text"`
}
type Track struct {
	Artist Artist  `json:"artist"`
	Images []Image `json:"image"`
	URL    string  `json:"url"`
	Name   string  `json:"name"`
}

func (t Track) ImageBySize(size string) string {
	for _, img := range t.Images {
		if img.Size == size {
			return img.URL
		}
	}
	return ""
}

type CurrentTrack struct {
	mu    sync.Mutex
	track TrackInfo
}

func (ct *CurrentTrack) Get() TrackInfo {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.track
}

func (ct *CurrentTrack) Set(t TrackInfo) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.track = t
}

type RecentTracks struct {
	Tracks []Track `json:"track"`
}

type LastFMResponse struct {
	RecentTracks RecentTracks `json:"recenttracks"`
}

func getScrobbles(user string, key string) ([]Track, error) {
	requestURL := fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=%s&api_key=%s&format=json", user, key)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not complete request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server returned error: %d %s", resp.StatusCode, resp.Body)
	}
	var result LastFMResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("could not decode request: %w", err)
	}
	return result.RecentTracks.Tracks, nil
}
