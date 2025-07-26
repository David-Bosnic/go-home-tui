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
	type event struct {
		Title     string `json:"title"`
		StartTime string `json:"startTime"`
		Location  string `json:"location"`
		EndTime   string `json:"endTime"`
	}
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", config.calendarID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("GET /calendar/events Error creating new req %v\n", err)
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
			return
		}
		if weeks <= 0 {
			log.Printf("GET /calendar/events Weeks is %v, weeks needs to be more then 0", weeks)
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
		return
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode == 401 {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(401)
		w.Write([]byte("Refresh OAuth Secret"))
		return
	}
	if res.StatusCode > 200 {
		log.Printf("GET /calendar/events Error failed with status code %v\n with body %v\n", res.StatusCode, body)
		return
	}
	if err != nil {
		log.Printf("GET /calendar/events Error reading body %v\n", err)
		return
	}
	var calendarEvent CalendarEvent
	err = json.Unmarshal(body, &calendarEvent)
	if err != nil {
		log.Printf("GET /calendar/events Error unmarshaling body %v\n", err)
		return
	}
	var events []event
	for _, item := range calendarEvent.Items {
		events = append(events, event{
			Title:     item.Summary,
			StartTime: item.Start.DateTime,
			EndTime:   item.End.DateTime,
			Location:  item.Location,
		})
	}
	dat, err := json.Marshal(events)
	if err != nil {
		log.Printf("Error marshaling chirps JSON: %s", err)
		w.Write(dat)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(dat)
}
