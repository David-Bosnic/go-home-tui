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
	apiConf.accessToken = "Bearer " + os.Getenv("ACCESS_TOKEN")
	apiConf.calendarID = os.Getenv("CALENDAR_ID")
	apiConf.refreshToken = os.Getenv("REFRESH_TOKEN")
	apiConf.clientID = os.Getenv("CLIENT_ID")
	apiConf.clientSecret = os.Getenv("CLIENT_SECRET")
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
