package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (config *apiConfig) handlerEventsGet(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GET /calendar/events Error creating new req %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", config.accessToken)

	queryParams := r.URL.Query()
	weeksParam := queryParams.Get("weeks")
	var weeks int
	var timeMax string
	if weeksParam != "" {
		weeks, err = strconv.Atoi(weeksParam)
		if err != nil {
			log.Printf("GET /calendar/events Error converting weeks to int %v\n", err)
			http.Error(w, "Invalid weeks parameter: must be a valid integer", http.StatusBadRequest)
			return
		}
		if weeks <= 0 {
			log.Printf("GET /calendar/events Weeks is %v, weeks needs to be more then 0", weeks)
			http.Error(w, "Invalid weeks parameter: must be greater than 0", http.StatusBadRequest)
			return
		}
		timeMax = time.Now().AddDate(0, 0, weeks*7).UTC().Format(time.RFC3339)
	} else {
		timeMax = time.Now().AddDate(0, 0, 7).UTC().Format(time.RFC3339)
	}

	q := req.URL.Query()
	q.Add("timeMin", time.Now().UTC().Format(time.RFC3339))
	q.Add("timeMax", timeMax)
	q.Add("orderBy", "startTime")
	q.Add("singleEvents", "true")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("GET /calendar/events Error fetching data %v\n", err)
		http.Error(w, "Failed to fetch calendar data", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("GET /calendar/events Error reading body %v\n", err)
		http.Error(w, "Failed to read calendar response", http.StatusInternalServerError)
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Refresh OAuth Secret"))
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("GET /calendar/events Error failed with status code %v\n with body %v\n", res.StatusCode, string(body))
		http.Error(w, "Calendar service error", http.StatusBadGateway)
		return
	}

	var calendarEvent CalendarEvent
	err = json.Unmarshal(body, &calendarEvent)
	if err != nil {
		log.Printf("GET /calendar/events Error unmarshaling body %v\n", err)
		http.Error(w, "Failed to parse calendar response", http.StatusInternalServerError)
		return
	}

	var events []event
	for _, item := range calendarEvent.Items {
		parsedTimeStart, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			log.Printf("GET /calendar/events Error parsing time %v\n", err)
			http.Error(w, "Failed to parse calendar time", http.StatusInternalServerError)
		}
		parsedTimeEnd, err := time.Parse(time.RFC3339, item.End.DateTime)
		if err != nil {
			log.Printf("GET /calendar/events Error parsing time %v\n", err)
			http.Error(w, "Failed to parse calendar time", http.StatusInternalServerError)
		}

		events = append(events, event{
			Title:     item.Summary,
			Date:      parsedTimeStart.Format("2006-01-02"),
			StartTime: parsedTimeStart.Format("15:04:05"),
			EndTime:   parsedTimeEnd.Format("15:04:05"),
			Location:  item.Location,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("GET /calendar/events Error encoding response JSON: %v\n", err)
		return
	}
}
