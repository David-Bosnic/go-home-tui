package main

import "os"

func loadConfig() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	url := configDir + "/go-home"
	err = os.MkdirAll(url, 0755)
	if err != nil {
		return "", err
	}
	dump := []byte("Here is less bytes")
	os.WriteFile(url+"/config", dump, 0755)
	return configDir, nil
}
