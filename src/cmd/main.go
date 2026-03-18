package main

import (
	"encoding/json"
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
	testFile := filepath.Join(wd, "examples", "config_example.txt")
	fileContent, err := os.ReadFile(testFile)
	if err != nil {
		log.Fatal(err)
	}
	status, err := src.ParseConfigStdOut(string(fileContent))
	if err != nil {
		log.Fatal(err)
	}
	out, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
