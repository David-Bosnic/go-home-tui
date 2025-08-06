package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Event struct {
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	EndTime   string `json:"endTime"`
}

type Point struct {
	x int
	y int
}

type Model struct {
	events []Event
	cursor Point
	point  Point

	selected map[Point]struct{}
}

var dayStyle = lipgloss.NewStyle().
	PaddingRight(5).
	PaddingLeft(4).
	Align(lipgloss.Center)

var addEventStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true, true, false, true).
	Width(10).
	Height(1).
	Align(lipgloss.Center)

var cardEventStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true, true, false, true).
	Align(lipgloss.Center).
	Width(10).
	Height(3)

var emptyEventStyle = lipgloss.NewStyle().
	PaddingRight(8).
	PaddingLeft(4).
	Align(lipgloss.Center)

var hoverAddEventStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("#6495ED")).
	Inherit(addEventStyle)

var hoverCardEventStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("#6495ED")).
	Inherit(cardEventStyle)

var hoverEmptyEventStyle = lipgloss.NewStyle().
	BorderForeground(lipgloss.Color("#6495ED")).
	Inherit(emptyEventStyle)

var hovered = lipgloss.NewStyle().
	Height(8).
	BorderBottom(true).
	BorderForeground(lipgloss.Color("#6495ED")).
	Inherit(cardEventStyle)

var whiteText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA"))

var redText = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000"))

func init() {
	SpinUp()
}

func initalModal(events []Event) Model {
	return Model{
		events:   events,
		selected: make(map[Point]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor.y > 0 {
				m.cursor.y--
			}

		case "down", "j":
			if m.cursor.y < len(m.events)-1 {
				m.cursor.y++
			}
		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
			}

		case "right", "l":
			if m.cursor.x < 7 {
				m.cursor.x++
			}

		case " ", "enter":
			_, ok := m.selected[Point{x: m.cursor.x, y: m.cursor.y}]
			if ok {
				delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
			} else {
				m.selected[Point{x: m.cursor.x, y: m.cursor.y}] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	s := whiteText.Render("Ye welp here are ye events")
	s += "\n\n"
	rows := eventRowCount(m.events)
	cols := 7
	eventMatrix := make([][]Event, rows)

	for i := range eventMatrix {
		eventMatrix[i] = make([]Event, cols)
	}

	dayMap := make(map[int]int)
	for _, event := range m.events {
		eventIndex := dateToIndex(event.Date)
		eventMatrix[dayMap[eventIndex]][eventIndex] = event
		dayMap[eventIndex]++
	}

	styledDays := getDaysStartingToday()
	for i := range styledDays {
		styledDays[i] = dayStyle.Render(styledDays[i])
	}
	s += lipgloss.JoinHorizontal(
		lipgloss.Top,
		styledDays...,
	)
	s += "\n"

	addEventCards := make([]Event, 7)
	for i := range addEventCards {
		addEventCards[i].Title = "+"
	}
	eventMatrix = append([][]Event{addEventCards}, eventMatrix...)
	for i, rows := range eventMatrix {
		rowEventsTitle := []string{}
		for j, event := range rows {
			currentPoint := Point{x: j, y: i}
			if m.cursor == currentPoint {
				switch event.Title {
				case "":
					rowEventsTitle = append(rowEventsTitle, hoverEmptyEventStyle.Render(""))
				case "+":
					rowEventsTitle = append(rowEventsTitle, hoverAddEventStyle.Render(event.Title))
				default:
					rowEventsTitle = append(rowEventsTitle, hoverCardEventStyle.Render(event.Title))
				}
			} else {
				switch event.Title {
				case "":
					rowEventsTitle = append(rowEventsTitle, emptyEventStyle.Render(""))
				case "+":
					rowEventsTitle = append(rowEventsTitle, addEventStyle.Render(event.Title))
				default:
					rowEventsTitle = append(rowEventsTitle, cardEventStyle.Render(event.Title))
				}

			}
		}
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			rowEventsTitle...,
		)
		s += "\n"
	}

	// if m.cursor == i {
	// 	if _, ok := m.selected[i]; ok {
	// 		s += hovered.Render(fmt.Sprintf("%s", event.Location))
	// 	} else {
	// 		s += hovered.Render(fmt.Sprintf("%s @ %s", event.Title, event.StartTime))
	// 	}
	// } else {
	// 	if _, ok := m.selected[i]; ok {
	// 		s += cardEventStyle.Render(fmt.Sprintf("%s", event.Location))
	// 	} else {
	// 		s += cardEventStyle.Render(fmt.Sprintf("%d %s", dateToIndex(event.Date), event.Title))
	// 	}
	// }
	// s += "\n"
	return s
}

func main() {
	refreshOauth()
	events := getEvents()
	p := tea.NewProgram(initalModal(events))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	// postEvent()

}
func getEvents() []Event {
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
func postEvent() {
	_, err := http.Post("http://localhost:8080/calendar/events", "", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func refreshOauth() {
	_, err := http.Post("http://localhost:8080/admin/refresh", "", nil)
	if err != nil {
		fmt.Println("Error with post req", err)
		os.Exit(1)
	}
}

func eventRowCount(events []Event) int {
	countMap := make(map[int]int)
	maxCount := 0
	for _, event := range events {
		countMap[dateToIndex(event.Date)]++
		if countMap[dateToIndex(event.Date)] > maxCount {
			maxCount++
		}
	}
	return maxCount
}
func getDaysStartingToday() []string {
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

func dateToIndex(date string) int {
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
