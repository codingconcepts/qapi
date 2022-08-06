package main

import (
	"flag"
	"log"
	"os"

	"github.com/codingconcepts/qapi/runner"
	"gopkg.in/yaml.v3"
)

var version string

func main() {
	log.SetFlags(0)

	configPath := flag.String("config", "", "absolution or relative path to the config file")
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

	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	var r runner.Runner
	if err := yaml.Unmarshal(configFile, &r); err != nil {
		log.Fatalf("error parsing config file: %v", err)
	}

	if err = r.Start(); err != nil {
		log.Fatalf("error executing runner: %v", err)
	}
}
