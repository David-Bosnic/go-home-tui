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
type model struct {
	events   []Event
	cursor   int
	selected map[int]struct{}
}

var cardEventStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder(), true, true, false, true).
	Width(20).
	Height(1)

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

func initalModal(events []Event) model {
	return model{
		events:   events,
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.events)-1 {
				m.cursor++
			}

		case " ", "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := whiteText.Render("Ye welp here are ye events")
	s += "\n\n"
	rows := eventRowCount(m.events)
	cols := 6
	matrix := make([][]Event, rows)
	for i := range matrix {
		matrix[i] = make([]Event, cols)
	}

	currRow := 0
	currCol := 0
	for _, event := range m.events {
		if dateToIndex(event.Date) == currRow || currRow == len(matrix) {
			matrix[currRow][currCol] = event
			currCol++
		} else {
			currRow++
			matrix[currRow][currCol] = event
			currCol++
		}
	}
	for _, row := range matrix {
		for _, value := range row {
			if value.Title == "" {
				continue
			}
			s += "hi"
		}
		s += "\n"
	}
	// 	if m.cursor == i {
	// 		if _, ok := m.selected[i]; ok {
	// 			s += hovered.Render(fmt.Sprintf("%s", event.Location))
	// 		} else {
	// 			s += hovered.Render(fmt.Sprintf("%s @ %s", event.Title, event.StartTime))
	// 		}
	// 	} else {
	// 		if _, ok := m.selected[i]; ok {
	// 			s += eventStyle.Render(fmt.Sprintf("%s", event.Location))
	// 		} else {
	// 			s += eventStyle.Render(fmt.Sprintf("%d %s", dateToIndex(event.Date), event.Title))
	// 		}
	// 	}
	// 	s += "\n"
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
