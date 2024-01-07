package system

import (
	"encoding/json"
	"io"
	"os"
)

// Write json to the given filename using efficient way to encode json data
func WriteJson[R any](filepath string, jsondata map[string]R) error {
	jsonString, err := json.MarshalIndent(jsondata, "", "  ")
	if err != nil {
		return err
	}

	writeErr := os.WriteFile(filepath, jsonString, os.ModePerm)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// Read json from the filepath
func ReadJson[R any](filepath string) (map[string]R, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var result map[string]R
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}
