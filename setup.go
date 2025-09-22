package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Setup() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		os.Exit(1)
	}
	url := configDir + "/go-home/.env"
	err = godotenv.Load(url)
	if err != nil {
		createConfig()
		fmt.Println("Creating Config in user config {USER_CONFIG}/go-home/.env")
		fmt.Println("Use README.md to config your credentials")
		os.Exit(1)
	}
	authFlag := flag.Bool("a", false, "Open Google Oauth on the Browser")
	flag.Parse()
	if *authFlag {
		fmt.Println("> Opening Google Oauth using default browser\n> http://localhost:8080/auth/google")
		OpenUrl("http://localhost:8080/auth/google")
		OauthSpinUp()
	} else {
		apiConf.accessToken = "Bearer " + os.Getenv("ACCESS_TOKEN")
		apiConf.calendarID = os.Getenv("CALENDAR_ID")
		apiConf.refreshToken = os.Getenv("REFRESH_TOKEN")
		apiConf.clientID = os.Getenv("CLIENT_ID")
		apiConf.clientSecret = os.Getenv("CLIENT_SECRET")
		err := RefreshOauth(apiConf)
		if err != nil {
			log.Printf("Failed to refresh Oauth %e", err)
			os.Exit(1)
		}
	}
}
