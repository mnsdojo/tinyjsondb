package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// Load reads a JSON file from the specified filePath and decodes it into the target map.
func Load(filepath string, target *map[string]interface{}) error {
	if filepath == "" {
		return errors.New("filepath cannot be empty")
	}

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func Save(filePath string, data map[string]interface{}) error {
	if filePath == "" {
		return errors.New("filePath cannot be empty")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	// Indentation for human readable formatting
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
