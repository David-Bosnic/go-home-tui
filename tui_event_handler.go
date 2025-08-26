package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func PostEvent(event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = http.Post("http://localhost:8080/calendar/events", "", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return nil
}

func DeleteEvent(event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	client := http.Client{}
	url := "http://localhost:8080/calendar/events"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func UpdateEvent(event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	client := &http.Client{}
	url := "http://localhost:8080/calendar/events"

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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
func GetEvents() []Event {
	res, err := http.Get("http://localhost:8080/calendar/events")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return events
}
