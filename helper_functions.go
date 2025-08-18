package main

import (
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
func FormsValidation(inputs []textinput.Model) error {
	startTime := inputs[StartTime].Value()
	endTime := inputs[EndTime].Value()
	_, err := time.Parse("15:04", startTime)
	if err != nil {
		return err
	}
	_, err = time.Parse("15:04", endTime)
	if err != nil {
		return err
	}
	return nil
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

func EventTimeFormatter(t string) int {
	return time.Time.Day(time.Now())
}
