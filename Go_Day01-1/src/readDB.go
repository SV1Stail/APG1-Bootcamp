package main

// package main

import (
	"day1/dbreader"
	"flag"
	"fmt"
	"os"
	"strings"
)

// type DBxml struct {
// }

// read path to file
func read_path() string {
	var f = flag.String("f", "", "read path")
	flag.Parse()
	return *f
}

// func (xml DBxml) DBRead(str string) (int, error) {

// }
func main() {
	path := read_path()
	if path == "" {
		fmt.Println("Please provide a file path using -f option")
		return
	}
	var db dbreader.Recipes
	var err error
	switch {
	case strings.HasSuffix(path, "json"):
		reader := dbreader.DB_json_reader{}
		db, err = reader.DBRead(path)
		reader.Converter_output(db)

	case strings.HasSuffix(path, "xml"):
		reader := dbreader.DB_xml_reader{}
		db, err = reader.DBRead(path)
		reader.Converter_output(db)
	default:
		fmt.Println("unsupported file type")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error reading database: %v\n", err)
		os.Exit(1)

	}
}
