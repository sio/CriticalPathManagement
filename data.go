package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Project struct {
	events     []Event
	activities []Activity
}

type Event struct {
	EarlyTime  int
	LatestTime int
}

type Activity struct {
	ID           string
	Description  string
	Dependencies []string
	Duration     int
}

func (a *Activity) Parse(row []string, offset offset) (err error) {
	if !offset.Valid() {
		return fmt.Errorf("not a valid offset: %v", offset)
	}
	a.ID = row[offset.ID]
	a.Description = row[offset.Description]
	a.Duration, err = strconv.Atoi(row[offset.Duration])
	if err != nil {
		return fmt.Errorf("can not parse activity duration (%s) for row %v: %w", row[offset.Duration], row, err)
	}
	a.Dependencies = make([]string, 0)
	ignore := map[string]bool{
		"-": true,
		"":  true,
	}
	for _, ID := range strings.Split(row[offset.Dependencies], ",") {
		ID = strings.TrimSpace(ID)
		if ignore[ID] {
			continue
		}
		a.Dependencies = append(a.Dependencies, ID)
	}
	return nil
}
