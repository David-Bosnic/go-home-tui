package main

import (
	"fmt"
	"testing"
)

func TestLoadingConfig(t *testing.T) {
	configDir, err := loadConfig()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("configDir: %v\n", configDir)
}
