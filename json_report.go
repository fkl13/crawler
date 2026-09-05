package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	if len(pages) == 0 {
		fmt.Println("No data to write to JSON")
		return nil
	}

	keys := slices.Collect(maps.Keys(pages))
	sort.Strings(keys)

	sorted := make([]PageData, 0, len(pages))
	for _, key := range keys {
		sorted = append(sorted, pages[key])
	}

	data, err := json.MarshalIndent(sorted, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = os.WriteFile(filename, data, 0o644)
	if err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	fmt.Printf("Report written to %s\n", filename)
	return nil
}
