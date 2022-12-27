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

	proj.Update()
	fmt.Printf("%v", proj)
}
