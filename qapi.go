package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codingconcepts/qapi/models"
	"github.com/codingconcepts/qapi/runner"
	"gopkg.in/yaml.v3"
)

var version string

func main() {
	log.SetFlags(0)

	configPath := flag.String("config", "", "absolute or relative path to a config file")
	showVersion := flag.Bool("version", false, "show the version")
	flag.Parse()

	if *showVersion {
		log.Println(version)
		os.Exit(0)
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := runFile(*configPath); err != nil {
		log.Fatalf("error running file: %v", err)
	}
}

func runFile(path string) error {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	var c models.Config
	if err := yaml.Unmarshal(configFile, &c); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	r := runner.New(&c)
	if err = r.Start(); err != nil {
		return fmt.Errorf("executing runner: %w", err)
	}

	return nil
}
