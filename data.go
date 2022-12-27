package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ID string

type Project struct {
	activities map[ID]*Activity
	events     []*Event
	start      *Event
	end        *Event
}

type Event struct {
	Index      int
	EarlyTime  int
	LatestTime int
}

type Activity struct {
	ID           ID
	Description  string
	Dependencies []ID
	Duration     int
	start        *Event
	end          *Event
}

func (a *Activity) Parse(row []string, offset offset) (err error) {
	if !offset.Valid() {
		return fmt.Errorf("not a valid offset: %v", offset)
	}
	a.ID = ID(row[offset.ID])
	a.Description = row[offset.Description]
	a.Duration, err = strconv.Atoi(row[offset.Duration])
	if err != nil {
		return fmt.Errorf("can not parse activity duration (%s) for row %v: %w", row[offset.Duration], row, err)
	}
	a.Dependencies = make([]ID, 0)
	ignore := map[string]bool{
		"-": true,
		"":  true,
	}
	for _, id := range strings.Split(row[offset.Dependencies], ",") {
		id = strings.TrimSpace(id)
		if ignore[id] {
			continue
		}
		a.Dependencies = append(a.Dependencies, ID(id))
	}
	return nil
}

func (p Project) String() string {
	return fmt.Sprintf("<Project: %d activities>", len(p.activities))
}

func (p *Project) Add(a *Activity) {
	if p.activities == nil {
		p.activities = make(map[ID]*Activity)
	}
	p.activities[a.ID] = a
}

func (p *Project) Update() error {
	const timeout = 2 * time.Second

	done := make(chan bool, 1)
	go func() {
		p.findEvents()
		p.findCriticalPath()
		done <- true
	}()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("updating timed out (%v)", timeout)
	}
}

func (p *Project) findEvents() {
	if p.start == nil {
		p.start = &Event{}
	}
	if p.end == nil {
		p.end = &Event{}
	}

	next := make(map[ID][]ID)
	for _, a := range p.activities {
		for _, dep := range a.Dependencies {
			next[dep] = append(next[dep], a.ID)
		}
	}

	var done bool
	for !done {
		for _, a := range p.activities {
			if a.start != nil && a.end != nil {
				continue
			}
			if len(a.Dependencies) == 0 {
				a.start = p.start
			}
			if len(next[a.ID]) == 0 {
				a.end = p.end
			}
			if a.end != nil && a.start == nil {
				a.start = &Event{}
				for _, id := range a.Dependencies {
					if p.activities[id].end != nil {
						a.start = p.activities[id].end
					}
					p.activities[id].end = a.start
				}
			}
		}
		done = true
		for _, a := range p.activities {
			if a.start == nil || a.end == nil {
				done = false
				break
			}
		}
	}

	p.straightenLooseEnds(p.end)

	seen := make(map[*Event]bool)
	for _, a := range p.activities {
		if !seen[a.start] {
			p.events = append(p.events, a.start)
			seen[a.start] = true
		}
		if !seen[a.end] {
			p.events = append(p.events, a.end)
			seen[a.end] = true
		}
	}
}

func (p *Project) straightenLooseEnds(cursor *Event) {
	for _, a := range p.activities {
		if a.end != cursor {
			continue
		}
		for _, id := range a.Dependencies {
			dep := p.activities[id]
			if dep.end != a.start {
				log.Printf("fixing a loose end on activity %s\n", dep.ID)
				dep.end = a.start
			}
		}
		if a.start != p.start {
			p.straightenLooseEnds(a.start)
		}
	}
}

func (p *Project) DebugPrint() {
	index := 1
	for _, a := range p.activities {
		if a.start.Index == 0 {
			a.start.Index = index
			index++
		}
		if a.end.Index == 0 {
			a.end.Index = index
			index++
		}
		fmt.Printf("Activity %s [%v->%v] len=%d\n", a.ID, a.start, a.end, a.Duration)
	}
}

func (p *Project) Critical(activity *Activity) bool {
	if activity.start.EarlyTime != activity.start.LatestTime {
		return false
	}
	if activity.end.EarlyTime != activity.end.LatestTime {
		return false
	}
	if activity.Duration != activity.end.EarlyTime-activity.start.EarlyTime {
		return false
	}
	return true
}

func (p *Project) findCriticalPath() {
	p.forwardPass(p.start)

	for _, a := range p.activities {
		a.end.LatestTime = p.end.EarlyTime
	}

	p.backwardPass(p.end)

	sort.Slice(p.events, func(i, j int) bool {
		return p.events[i].EarlyTime < p.events[j].EarlyTime
	})
	for index, event := range p.events {
		event.Index = index + 1
	}
}

func (p *Project) forwardPass(cursor *Event) {
	for _, a := range p.activities {
		if a.start != cursor {
			continue
		}
		value := a.start.EarlyTime + a.Duration
		if value > a.end.EarlyTime {
			a.end.EarlyTime = value
		}
		p.forwardPass(a.end)
	}
}

func (p *Project) backwardPass(cursor *Event) {
	for _, a := range p.activities {
		if a.end != cursor {
			continue
		}
		value := a.end.LatestTime - a.Duration
		if value < a.start.LatestTime {
			a.start.LatestTime = value
		}
		p.backwardPass(a.start)
	}
}
