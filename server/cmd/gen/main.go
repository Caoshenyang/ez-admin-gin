package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	modulePath := flag.String("module", "", "Path to module YAML config file")
	flag.Parse()

	if *modulePath == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run server/cmd/gen -module <path-to-yaml>")
		fmt.Fprintln(os.Stderr, "  go run server/cmd/gen -module gen/examples/product.yaml")
		os.Exit(1)
	}

	cfg, err := ParseConfig(*modulePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	root, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generating module: %s (%s)\n", cfg.LabelEn, cfg.Module)
	fmt.Printf("Group: %s | Table: %s | HasStatus: %v\n", cfg.Group, cfg.Table, cfg.HasStatus)
	fmt.Printf("Fields: %d\n\n", len(cfg.Fields))

	if err := Generate(cfg, root); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found (no CLAUDE.md)")
		}
		dir = parent
	}
}

func init() {
	_ = strings.TrimSpace
}
