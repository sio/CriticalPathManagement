package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const debug = false

func main() {
	defaultSource := strings.TrimSuffix(os.Args[0], filepath.Ext(os.Args[0])) + ".xlsx"
	filename := flag.String("f", defaultSource, "path to Excel file with project data")
	flag.Parse()

	var err error
	if debug {
		logFile, err := os.OpenFile(os.Args[0]+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	_, err = os.Stat(*filename)
	if os.IsNotExist(err) {
		log.Printf("file does not exist: %s, trying with ", *filename)
		*filename = filepath.Join(filepath.Dir(os.Args[0]), filepath.Base(*filename))
		log.Printf("%s\n", *filename)
	}
	log.Printf("opening Excel file: %s\n", *filename)

	var proj Project

	proj, err = ReadExcel(*filename)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = proj.Update()
	if err != nil {
		log.Fatalf("%v", err)
	}

	dest := strings.TrimSuffix(*filename, filepath.Ext(*filename)) + ".svg"
	log.Printf("rendering graph to %s\n", dest)
	err = proj.Render(dest)
	if err != nil {
		log.Fatal(err)
	}

	proj.DebugPrint()
	log.Printf("minimum project duration: %d", proj.end.EarlyTime)
}
