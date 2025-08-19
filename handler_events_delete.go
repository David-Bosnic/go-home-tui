package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (config *apiConfig) handlerEventsDelete(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Printf("DELETE /calendar/events Error Decoding event %v\n", err)
		http.Error(w, "Failed to decode calendar event", http.StatusBadRequest)
		return
	}

	client := http.Client{}
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events/%s", config.calendarID, event.Id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", config.accessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("PATCH /calendar/events Error failed with status code %v\n", resp.StatusCode)
		http.Error(w, "Calendar service error", http.StatusBadGateway)
		return
	}
}
