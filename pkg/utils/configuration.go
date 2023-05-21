package utils

import (
	"encoding/json"
	"os"
)

func LoadSettings(file string, settings interface{}, defaultSettings interface{}) error {
	var err error
	if !FileExists(file) {
		err = CreateSettings(file, defaultSettings)
	}
	if err != nil {
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, settings)
	if err != nil {
		return err
	}

	return nil
}

func CreateSettings(file string, defaultSettings interface{}) error {
	content, err := json.MarshalIndent(defaultSettings, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
