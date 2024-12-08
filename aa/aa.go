package aa

import (
	"fmt"
	"regexp"
	"time"

	"github.com/xackery/critsprinkler/tracker"
)

type AA struct {
	lastZoneEvent time.Time
	lastAAExpDing time.Time
	parseStart    time.Time
	totalAAGained int
	zone          string
}

var (
	instance    *AA
	aaGainRegex = regexp.MustCompile(`You have gained an ability point!`)
)

func New() (*AA, error) {
	if instance != nil {
		return nil, fmt.Errorf("aa already exists")
	}
	a := &AA{
		zone:          "Unknown",
		lastZoneEvent: time.Now(),
		parseStart:    time.Now(),
	}

	/* w, err := os.Create("aa.txt")
	if err != nil {
		return nil, fmt.Errorf("aa create file: %w", err)
	}
	defer w.Close()

	w.WriteString("") */

	err := tracker.Subscribe(a.onLine)
	if err != nil {
		return nil, fmt.Errorf("tracker subscribe: %w", err)
	}
	err = tracker.SubscribeToZoneEvent(a.onZone)
	if err != nil {
		return nil, fmt.Errorf("tracker subscribe to zone: %w", err)
	}
	instance = a
	return a, nil
}

func (a *AA) onLine(event time.Time, line string) {
	a.onAA(event, line)
}

func (a *AA) onZone(event time.Time, zoneName string) {
	a.lastZoneEvent = time.Now()
	a.zone = zoneName
	a.totalAAGained = 0
}

func (a *AA) onAA(event time.Time, line string) {
	match := aaGainRegex.FindStringSubmatch(line)
	if len(match) < 1 {
		return
	}

	event = time.Now()
	sinceLastDing := event.Sub(a.lastAAExpDing).Minutes()
	a.lastAAExpDing = event
	a.totalAAGained++

	elapsedTime := event.Sub(a.lastZoneEvent).Hours()
	if elapsedTime <= 0 {
		fmt.Println("Last zone event was", elapsedTime, "hours ago", a.lastZoneEvent, event)
		fmt.Println("Elapsed time is 0")
		return
	}
	aaPerHour := float64(a.totalAAGained) / elapsedTime
	fmt.Printf("Total AA gained: %d / per hour: %.2f, Time in zone: %0.2f hours, last AA ding: %0.2f minutes\n", a.totalAAGained, aaPerHour, elapsedTime, sinceLastDing)

	/* 	w, err := os.Create("aa.txt")
	   	if err != nil {
	   		fmt.Println(err)
	   		return
	   	}
	   	defer w.Close()

	   	w.WriteString(fmt.Sprintf("%d AA per hour in %0.2f hours\n", int(aaPerHour), elapsedTime)) */
}
