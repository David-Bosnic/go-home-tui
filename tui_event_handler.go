package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func PostEvent() {
	_, err := http.Post("http://localhost:8080/calendar/events", "", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func UpdateEvent(event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	url := "http://localhost:8080/calendar/events"

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

}
func RefreshOauth() {
	_, err := http.Post("http://localhost:8080/admin/refresh", "", nil)
	if err != nil {
		fmt.Println("Error with post req", err)
		os.Exit(1)
	}
}
