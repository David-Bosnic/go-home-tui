package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Event struct {
	Id        string `json:"eventId"`
	Summary   string `json:"summary"`
	StartTime string `json:"start_time"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	EndTime   string `json:"end_time"`
}

type Point struct {
	x int
	y int
}

type Model struct {
	events      []Event
	cursor      Point
	point       Point
	selected    map[Point]struct{}
	eventMatrix [][]Event
	mode        string
	inputs      []textinput.Model
	focusIndex  int
	flipState   bool
}

// Styles
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedButton = focusedStyle.Render("[ Submit ]")

	dayStyle = lipgloss.NewStyle().
			PaddingRight(7).
			PaddingLeft(7).
			Align(lipgloss.Center)

	addEventStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, false, true).
			Width(15).
			Height(1).
			Align(lipgloss.Center)

	cardEventStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, false, true).
			Align(lipgloss.Center).
			Width(15).
			Height(5)

	emptyEventStyle = lipgloss.NewStyle().
			PaddingRight(8).
			PaddingLeft(9).
			Align(lipgloss.Center)

	hoverAddEventStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#6495ED")).
				Inherit(addEventStyle)

	hoverCardEventStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#6495ED")).
				Inherit(cardEventStyle)

	hoverEmptyEventStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#6495ED")).
				Inherit(emptyEventStyle)

	hovered = lipgloss.NewStyle().
		Height(8).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#6495ED")).
		Inherit(cardEventStyle)

	whiteText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	redText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))
)

func init() {
	SpinUp()
}

func initialModel(events []Event) Model {

	rows := eventRowCount(events)
	cols := 7
	eventMatrix := make([][]Event, rows)

	for i := range eventMatrix {
		eventMatrix[i] = make([]Event, cols)
	}

	dayMap := make(map[int]int)
	for _, event := range events {
		eventIndex := dateToIndex(event.Date)
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

	m := Model{
		events:      events,
		selected:    make(map[Point]struct{}),
		eventMatrix: eventMatrix,
		mode:        "calendar",
		inputs:      make([]textinput.Model, 5),
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		m.inputs[i] = t
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.mode == "calendar" {
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
				if m.cursor.y < eventRowCount(m.events) && m.eventMatrix[m.cursor.y+1][m.cursor.x].Summary != "" {
					m.cursor.y++
				}
			case "left", "h":
				if m.cursor.x > 0 && m.eventMatrix[m.cursor.y][m.cursor.x-1].Summary != "" {
					m.cursor.x--
				}

			case "right", "l":
				if m.cursor.x < 6 && m.eventMatrix[m.cursor.y][m.cursor.x+1].Summary != "" {
					m.cursor.x++
				}

			case "f":
				m.flipState = !m.flipState

			case " ", "enter":
				_, ok := m.selected[Point{x: m.cursor.x, y: m.cursor.y}]
				if ok {
					delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
				} else {
					m.mode = "forms"
					m.focusIndex = 0
					event := m.eventMatrix[m.cursor.y][m.cursor.x]
					m.inputs[0].SetValue(event.Summary)
					m.inputs[1].SetValue(event.StartTime)
					m.inputs[2].SetValue(event.EndTime)
					m.inputs[3].SetValue(event.Location)
					m.inputs[4].SetValue(event.Id)
					m.selected[Point{x: m.cursor.x, y: m.cursor.y}] = struct{}{}
				}
			}
		}
	}
	if m.mode == "forms" {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()
				if s == "enter" && m.focusIndex == len(m.inputs) {
					UpdateEvent(Event{
						Id:      m.inputs[4].Value(),
						Summary: m.inputs[0].Value(),
					})
					m.mode = "calendar"
					return m, nil
				}

				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				}

				if s == "down" {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs)
				}

				cmds := make([]tea.Cmd, len(m.inputs))
				for i := 0; i <= len(m.inputs)-1; i++ {
					if i == m.focusIndex {
						cmds[i] = m.inputs[i].Focus()
						m.inputs[i].PromptStyle = focusedStyle
						m.inputs[i].TextStyle = focusedStyle
						continue
					}
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle

				}
				return m, tea.Batch(cmds...)
			}
			cmd := m.updateInputs(msg)
			return m, cmd
		}

	}
	return m, nil

}
func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string
	if m.mode == "forms" {

		labels := []string{"Event:", "Start Time:", "End Time:", "Location:", "Id: "}
		for i := range m.inputs {
			s += labels[i] + "\n" + m.inputs[i].View()
			if i < len(m.inputs)-1 {
				s += "\n"
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		var b strings.Builder
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)
		s += b.String()
	} else {

		s += whiteText.Render("Current Event:", m.eventMatrix[m.cursor.y][m.cursor.x].Summary)
		s += "\n\n"

		styledDays := getDaysStartingToday()
		for i := range styledDays {
			styledDays[i] = dayStyle.Render(styledDays[i])
		}
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			styledDays...,
		)
		s += "\n"

		for i, rows := range m.eventMatrix {
			rowEventsTitle := []string{}
			for j, event := range rows {
				currentPoint := Point{x: j, y: i}
				if m.cursor == currentPoint {
					switch event.Summary {
					case "":
						rowEventsTitle = append(rowEventsTitle, hoverEmptyEventStyle.Render(""))
					case "+":
						rowEventsTitle = append(rowEventsTitle, hoverAddEventStyle.Render(event.Summary))
					default:
						//TODO: Maybe truncate super long event names
						if !m.flipState {
							rowEventsTitle = append(rowEventsTitle, hoverCardEventStyle.Render(truncate(event.Summary, 35, false), event.StartTime+"-"+event.EndTime))
						} else {
							rowEventsTitle = append(rowEventsTitle, hoverCardEventStyle.Render(event.Location))
						}
					}
					continue
				} else {
					switch event.Summary {
					case "":
						rowEventsTitle = append(rowEventsTitle, emptyEventStyle.Render(""))
					case "+":
						rowEventsTitle = append(rowEventsTitle, addEventStyle.Render((event.Summary)))
					default:
						rowEventsTitle = append(rowEventsTitle, cardEventStyle.Render(truncate(event.Summary, 26, true)))
					}

				}
			}
			s += lipgloss.JoinHorizontal(
				lipgloss.Top,
				rowEventsTitle...,
			)
			s += "\n"
		}
	}
	return s
}

func main() {
	RefreshOauth()
	events := getEvents()
	p := tea.NewProgram(initialModel(events))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
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

func truncate(s string, maxLen int, elipse bool) string {
	if len(s) <= maxLen {
		return s
	}
	if elipse {
		return s[:maxLen-3] + "\n..."
	}
	return s[:maxLen]
}
