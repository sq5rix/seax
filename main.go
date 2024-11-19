package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/seax/client"
)

const (
	defaultInstanceURL = "http://localhost:4000"
	defaultTimeout    = 10 * time.Second
)

type config struct {
	instanceURL string
	format     string
	timeout    time.Duration
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := parseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Read the query from stdin
	query, err := readQuery()
	if err != nil {
		return err
	}

	// Initialize the SearXNG client with timeout
	searxClient := client.NewClient(cfg.instanceURL, client.WithTimeout(cfg.timeout))

	// Perform the search
	results, err := searxClient.Search(query)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Output results based on format
	return outputResults(results, cfg.format)
}

func parseFlags() (config, error) {
	cfg := config{}
	
	flag.StringVar(&cfg.instanceURL, "url", defaultInstanceURL, "SearXNG instance URL")
	flag.StringVar(&cfg.format, "format", "json", "Output format (json or text)")
	timeout := flag.Duration("timeout", defaultTimeout, "Search timeout duration")
	
	flag.Parse()
	
	cfg.timeout = *timeout

	// Validate format
	if cfg.format != "json" && cfg.format != "text" {
		return cfg, fmt.Errorf("invalid format: %s (must be 'json' or 'text')", cfg.format)
	}

	return cfg, nil
}

func readQuery() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("search query cannot be empty")
	}
	
	return query, nil
}

func outputResults(results interface{}, format string) error {
	switch format {
	case "json":
		return outputJSON(results)
	case "text":
		return outputText(results)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func outputJSON(results interface{}) error {
	jsonOutput, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results to JSON: %w", err)
	}

	fmt.Println(string(jsonOutput))
	return nil
}

func outputText(results interface{}) error {
	// Type assert to expected result type and format as needed
	// This is a basic implementation - adjust according to your actual results structure
	fmt.Printf("%+v\n", results)
	return nil
}
