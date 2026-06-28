package main

type TrackInfo struct {
	Artist   string
	Title    string
	ImageURL string
}

type ServerStats struct {
	CPU    string
	Memory string
	Uptime string
}

func GetTrackInfo(t Track) TrackInfo {
	info := TrackInfo{
		Title:    t.Name,
		Artist:   t.Artist.Name,
		ImageURL: t.ImageBySize("large"),
	}
	return info
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[0:n] + "..."
	}
	return s
}
