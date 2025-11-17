package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codingconcepts/qapi/models"
	"github.com/codingconcepts/qapi/runner"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

var version string

func main() {
	log.SetFlags(0)

	configPath := flag.String("config", "", "absolute or relative path to a config file")
	vus := flag.Int("vus", 1, "the number of virtual users to simulate")
	duration := flag.Duration("duration", time.Minute*1, "length of test")
	showVersion := flag.Bool("version", false, "show the version")
	debug := flag.Bool("debug", false, "verbose logging for debugging")
	flag.Parse()

	if *showVersion {
		log.Println(version)
		os.Exit(0)
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsExclude: []string{
			zerolog.TimestampFieldName,
		},
	}).Level(lo.Ternary(*debug, zerolog.DebugLevel, zerolog.InfoLevel))

	config, err := parseConfig(*configPath)
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	events := make(chan models.RequestResult, *vus)

	printer := runner.NewPrinter(events, &logger)

	go printer.Run()

	if err := run(config, *vus, *duration, events, &logger); err != nil {
		log.Fatalf("error running: %v", err)
	}
}

func parseConfig(path string) (models.Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return models.Config{}, fmt.Errorf("reading config file: %w", err)
	}

	var c models.Config
	if err := yaml.Unmarshal(configFile, &c); err != nil {
		return models.Config{}, fmt.Errorf("unmarshalling yaml: %w", err)
	}

	return c, nil
}

func run(c models.Config, vus int, d time.Duration, e chan models.RequestResult, logger *zerolog.Logger) error {
	var eg errgroup.Group

	for range vus {
		eg.Go(func() error {
			r := runner.New(c, d, e, logger)
			if err := r.Start(); err != nil {
				return fmt.Errorf("executing runner: %w", err)
			}

			return nil
		})
	}

	return eg.Wait()
}
