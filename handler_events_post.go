package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PostEvent struct {
	Summary string `json:"summary"`
	Start   struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone,omitempty"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone,omitempty"`
	} `json:"end"`
}

func (config *apiConfig) handlerEventsPost(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)
	newEvent := PostEvent{
		Summary: "Test Event",
	}
	newEvent.Start.DateTime = "2025-08-01T10:00:00-06:00"

	newEvent.End.DateTime = "2025-08-01T11:00:00-06:00"

	out, err := json.Marshal(newEvent)
	if err != nil {
		log.Printf("POST /calendar/events Error marshaling event %v\n", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Printf("POST /calendar/events Error creating new req %v\n", err)
		return
	}
	fmt.Printf("\nHere is the payload %v\n\n", req)
	req.Header.Set("Authorization", config.accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GET /calendar/events Error fetching data %v\n", err)
		return
	}
	fmt.Println(res)
}
