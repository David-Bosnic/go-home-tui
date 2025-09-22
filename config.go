package main

import (
	"os"
	"path/filepath"
)

func createConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "go-home")
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		return err
	}
	dump := []byte(`ACCESS_TOKEN="ACCESS_TOKEN"
CALENDAR_ID="email@gmail.com"
CLIENT_ID="CLIENT_ID"
CLIENT_SECRET="CLIENT_SECRET"
REFRESH_TOKEN="REFRESH_TOKEN"
COLOR_PRIMARY="#7e9cd8"
COLOR_WARNING="#ffcc00"
COLOR_ERROR="#FF3333"
`)
	envPath := filepath.Join(configPath, ".env")
	os.WriteFile(envPath, dump, 0644)
	return nil
}
