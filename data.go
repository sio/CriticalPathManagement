package main

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
	Dependencies []Activity
	Duration     int
}
