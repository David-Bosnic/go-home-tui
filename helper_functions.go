package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

func GetDaysStartingToday() []string {
	allDays := []string{
		"Sun",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}
	today := int(time.Now().Weekday())
	return append(allDays[today:], allDays[:today]...)
}
func CreateEventMatrix(events []Event) [][]Event {
	rows := EventRowCount(events)
	cols := 7
	eventMatrix := make([][]Event, rows)

	for i := range eventMatrix {
		eventMatrix[i] = make([]Event, cols)
	}

	dayMap := make(map[int]int)
	for _, event := range events {
		eventIndex := DateToIndex(event.Start.Date)
		if eventIndex >= 0 && eventIndex < cols && dayMap[eventIndex] < rows {
			eventMatrix[dayMap[eventIndex]][eventIndex] = event
			dayMap[eventIndex]++
		}
	}

	addEventCards := make([]Event, 7)
	for i := range addEventCards {
		addEventCards[i].Summary = "+"
	}
	eventMatrix = append([][]Event{addEventCards}, eventMatrix...)
	return eventMatrix
}
func EventRowCount(events []Event) int {
	countMap := make(map[int]int)
	maxCount := 0
	for _, event := range events {
		countMap[DateToIndex(event.Start.DateTime.Format("2006-01-02"))]++
		if countMap[DateToIndex(event.Start.DateTime.Format("2006-01-02"))] > maxCount {
			maxCount++
		}
	}
	return maxCount
}
func DateToIndex(date string) int {
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return -1
	}

	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, now.Location())

	diff := targetDate.Sub(currentDate)
	days := int(diff.Hours() / 24)

	if days < 0 || days > 7 {
		return -1
	}

	return days
}
func FormsValidation(inputs []textinput.Model, validFields *[]bool) bool {
	for i := range *validFields {
		(*validFields)[i] = true
	}
	var invalid bool
	summary := inputs[Summary].Value()
	if summary == "" {
		(*validFields)[Summary] = false
		invalid = true
	}
	date := inputs[Date].Value()
	startTime := inputs[StartTime].Value()
	endTime := inputs[EndTime].Value()

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		(*validFields)[Date] = false
		invalid = true
	}
	_, err = time.Parse("15:04", startTime)
	if err != nil {
		(*validFields)[StartTime] = false
		invalid = true
	}
	_, err = time.Parse("15:04", endTime)
	if err != nil {
		(*validFields)[EndTime] = false
		invalid = true
	}
	return invalid
}
func Truncate(s string, maxLen int, elipse bool) string {
	if len(s) <= maxLen {
		return s
	}
	if elipse {
		return s[:maxLen-3] + "\n..."
	}
	return s[:maxLen]
}
func NewEventDate(i int) string {
	now := time.Now()
	eventDate := now.AddDate(0, 0, i)
	return eventDate.Format("2006-01-02")
}
func OpenUrl(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Unsupported Browser")
	}
	if err != nil {
		return err
	}
	return nil
}
