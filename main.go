package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	defaultSource := strings.TrimSuffix(os.Args[0], filepath.Ext(os.Args[0])) + ".xlsx"
	filename := flag.String("f", defaultSource, "path to Excel file with project data")
	flag.Parse()

	log.Printf("opening Excel file: %s\n", *filename)

	var err error
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

	fmt.Printf("%v\n", proj)
	proj.DebugPrint()
}
