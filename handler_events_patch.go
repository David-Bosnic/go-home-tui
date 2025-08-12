package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (config *apiConfig) handlerEventsPatch(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Printf("PATCH /calendar/events Error Decoding event %v\n", err)
		http.Error(w, "Failed to decode calendar event", http.StatusBadRequest)
		return
	}
	payload, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events/%s", config.calendarID, event.Id)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", config.accessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("PATCH /calendar/events Error failed with status code %v\n with body %v\n", resp.StatusCode, string(body))
		http.Error(w, "Calendar service error", http.StatusBadGateway)
		return
	}
	log.Println("Status code", resp.StatusCode, string(body))

}
