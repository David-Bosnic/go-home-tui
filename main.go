package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	Setup()
	godotenv.Load()
	style = SetStyles()
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
