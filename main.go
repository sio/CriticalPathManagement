package main

import (
	"flag"
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
	err = ReadExcel(*filename)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
