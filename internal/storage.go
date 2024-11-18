package storage

import (
	"encoding/json"
	"os"
)

// Load reads a JSON file from the specified filePath and decodes it into the target map.
func Load(filepath string, target *map[string]interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func Save(filePath string, data map[string]interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	// Indentation for human readable formatting
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}
