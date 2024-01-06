package system

import (
	"encoding/json"
	"os"
)

// Write json to the given filename using efficient way to encode json data
func WriteJson(filepath string, dataList []map[string]interface{}) error {
	file, err := os.OpenFile(filepath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encodeErr := encoder.Encode(dataList)

	if encodeErr != nil {
		return encodeErr
	}

	return nil
}

// Read json from the filepath
func ReadJson(filepath string) ([]map[string]interface{}, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	dataList := []map[string]interface{}{}

	// Read the array open bracket
	decoder.Token()

	data := map[string]interface{}{}
	for decoder.More() {
		decoder.Decode(&data)
		dataList = append(dataList, data)
	}

	return dataList, nil
}
