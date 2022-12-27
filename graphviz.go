package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

const criticalPathColor = "red"

func (p *Project) Render(path string) (err error) {
	gv := graphviz.New()
	graph, err := gv.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err = graph.Close(); err != nil {
			log.Fatal(err)
		}
		gv.Close()
	}()

	graph.SetRankDir(cgraph.LRRank)

	nodes := make(map[*Event]*cgraph.Node)
	for _, event := range p.events {
		nodes[event], err = graph.CreateNode(strconv.Itoa(event.Index))
		if err != nil {
			return fmt.Errorf("could not add node: %w", err)
		}
		if event.EarlyTime == event.LatestTime {
			nodes[event].SetColor(criticalPathColor)
			nodes[event].SetFontColor(criticalPathColor)
		}
		nodes[event].SetShape(cgraph.CircleShape)
	}
	for _, activity := range p.activities {
		edge, err := graph.CreateEdge(string(activity.ID), nodes[activity.start], nodes[activity.end])
		if err != nil {
			return fmt.Errorf("could not add edge: %w", err)
		}
		edge.SetLabel(string(activity.ID))
		//edge.SetMinLen(activity.Duration)
		if p.Critical(activity) {
			edge.SetColor(criticalPathColor)
			edge.SetFontColor(criticalPathColor)
		}
	}

	var buf bytes.Buffer
	err = gv.Render(graph, "dot", &buf)
	if err != nil {
		return fmt.Errorf("failed to render to dot: %w", err)
	}
	fmt.Println(buf.String())

	err = gv.RenderFilename(graph, graphviz.SVG, path)
	if err != nil {
		return fmt.Errorf("failed to render: %w", err)
	}
	return nil
}
