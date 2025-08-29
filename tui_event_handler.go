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

func PostEvent(event Event, config apiConfig) error {
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)

	var postEvent PostEventType
	postEvent.Summary = event.Summary
	postEvent.Location = event.Location
	postEvent.Start.DateTime = event.Start.DateTime.Format(time.RFC3339)
	postEvent.End.DateTime = event.End.DateTime.Format(time.RFC3339)

	payload, err := json.Marshal(postEvent)
	if err != nil {
		log.Printf("POST /calendar/events Error marshaling event %v\n", err)
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("POST /calendar/events Error creating new req %v\n", err)
		return err
	}
	req.Header.Set("Authorization", config.accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("POST /calendar/events Error making request %v\n", err)
		return err
	}
	defer res.Body.Close()
	return nil
}

func GetEvents(config apiConfig) ([]Event, error) {
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GET /calendar/events Error creating new req %v\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", config.accessToken)

	var timeMax string
	timeMax = time.Now().AddDate(0, 0, 7).UTC().Format(time.RFC3339)

	q := req.URL.Query()
	q.Add("timeMin", time.Now().UTC().Format(time.RFC3339))
	q.Add("timeMax", timeMax)
	q.Add("orderBy", "startTime")
	q.Add("singleEvents", "true")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GET /calendar/events Error fetching data %v\n", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("GET /calendar/events Error reading body %v\n", err)
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Printf("GET /calendar/events Error status Unauthorized %v\n", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("GET /calendar/events Error failed with status code %v\n with body %v\n", res.StatusCode, string(body))
		return nil, err
	}

	var calendarEvent CalendarEvent
	err = json.Unmarshal(body, &calendarEvent)
	if err != nil {
		log.Printf("GET /calendar/events Error unmarshaling body %v\n", err)
		return nil, err
	}

	var events []Event
	for _, item := range calendarEvent.Items {
		parsedTimeStart, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			log.Printf("GET /calendar/events Error parsing time %v\n", err)
			return nil, err
		}
		parsedTimeEnd, err := time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			log.Printf("GET /calendar/events Error parsing time %v\n", err)
			return nil, err
		}
		_, startZone := parsedTimeStart.Zone()
		_, endZone := parsedTimeEnd.Zone()
		events = append(events, Event{
			Id:      item.ID,
			Summary: item.Summary,
			Start: DateTime{
				DateTime: parsedTimeStart,
				Date:     parsedTimeStart.Format(time.DateOnly),
				TimeZone: startZone,
			},
			End: DateTime{
				DateTime: parsedTimeEnd,
				Date:     parsedTimeEnd.Format(time.DateOnly),
				TimeZone: endZone,
			},
			Location: item.Location,
		})
	}

	return events, nil
}
func DeleteEvent(event Event, config apiConfig) error {
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

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("DELETE /calendar/events Error failed with status code %v\n", resp.StatusCode)
		return err
	}
	return nil
}

func UpdateEvent(event Event, config apiConfig) error {
	var patchEvent PatchEventType
	patchEvent.Summary = event.Summary
	patchEvent.Location = event.Location
	patchEvent.Start.DateTime = event.Start.DateTime.Format(time.RFC3339)
	patchEvent.End.DateTime = event.End.DateTime.Format(time.RFC3339)

	payload, err := json.Marshal(patchEvent)
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
		return err
	}

	return nil
}

func RefreshOauth() error {
	resp, err := http.Post("http://localhost:8080/admin/refresh", "", nil)
	if err != nil {
		return fmt.Errorf("Failed to connect to local server")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed with status code %d", resp.StatusCode)
	}
	return nil
}
