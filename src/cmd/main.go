package main

import (
	"fmt"
	src "go_pwrstat_api/src"
	"log"
	"os"
	"path/filepath"
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	testFile := filepath.Join(wd, "examples", "status_example.txt")
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		log.Fatal(err)
	}
	status, err := src.ParseStatusStdOut(string(fileContent))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("status --> %+v\n", status)
}
