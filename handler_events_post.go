package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type PostEventType struct {
	Summary  string `json:"summary"`
	Location string `json:"location,omitempty"`
	Start    struct {
		DateTime string `json:"dateTime"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
	} `json:"end"`
}

func (config *apiConfig) handlerEventsPost(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)

	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Printf("POST /calendar/events Error Decoding event %v\n", err)
		http.Error(w, "Failed to decode calendar event", http.StatusBadRequest)
		return
	}
	var postEvent PostEventType
	postEvent.Summary = event.Summary
	postEvent.Location = event.Location
	postEvent.Start.DateTime = event.Start.DateTime.Format(time.RFC3339)
	postEvent.End.DateTime = event.End.DateTime.Format(time.RFC3339)

	payload, err := json.Marshal(postEvent)
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
