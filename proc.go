package main

import (
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type CurrentStats struct {
	mu    sync.Mutex
	stats ServerStats
}

func (cs *CurrentStats) Get() ServerStats {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.stats
}

func (cs *CurrentStats) Set(s ServerStats) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.stats = s
}

func getMemory() (string, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f/%.1f GB", float32(mem.Used)/1024/1024/1024, float32(mem.Total)/1024/1024/1024), nil
}
func getCPU() (string, error) {
	percent, err := cpu.Percent(0, false)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f%%", percent[0]), nil
}
func getUptime() (string, error) {
	uptimeSec, err := host.Uptime()
	if err != nil {
		return "", err
	}
	uptimeDays := float32(uptimeSec) / 3600 / 24
	// uptimeMin := uptimeSec / 60
	if uptimeDays == 1 {
		return fmt.Sprintf("%.1f day", uptimeDays), nil
	}
	return fmt.Sprintf("%.1f days", uptimeDays), nil
}
