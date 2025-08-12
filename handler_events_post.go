package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	payload, err := json.Marshal(newEvent)
	if err != nil {
		log.Printf("POST /calendar/events Error marshaling event %v\n", err)
		http.Error(w, "Failed to parse calendar event", http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("POST /calendar/events Error creating new req %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", config.accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("POST /calendar/events Error making request %v\n", err)
		http.Error(w, "Failed to post calendar data", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
