package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"strings"
	"time"
)

type DateTime struct {
	DateTime time.Time `json:"dateTime"`
	Date     string    `json:"date"`
	TimeZone int       `json:"timeZone"`
}
type Event struct {
	Id       string   `json:"event_id"`
	Summary  string   `json:"summary"`
	Start    DateTime `json:"start"`
	End      DateTime `json:"end"`
	Location string   `json:"location"`
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Flip  key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Flip: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "toggle location"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Point struct {
	x int
	y int
}

type Model struct {
	spinner      spinner.Model
	events       []Event
	keys         keyMap
	help         help.Model
	cursor       Point
	point        Point
	selected     map[Point]struct{}
	eventMatrix  [][]Event
	mode         string
	inputs       []textinput.Model
	focusIndex   int
	showLocation bool
	newEvent     bool
	validFields  []bool
	areYouSure   bool
}

const (
	Summary = iota
	Date
	StartTime
	EndTime
	Location
	Id
)

// Styles
var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7e9cd8"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#7e9cd8"))
	blurredSubmitButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	focusedSubmitButton = focusedStyle.Render("[ Submit ]")

	blurredCancelButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Cancel"))
	focusedCancelButton = focusedStyle.Render("[ Cancel ]")

	blurredDeleteButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Delete"))
	focusedDeleteButton = focusedStyle.Render("[ Delete ]")

	grayBlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	grayBlurredDeleteButton = grayBlurredStyle.Render("[ Delete ]")

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
			Width(15).
			Height(5).
			Align(lipgloss.Center)

	emptyEventStyle = lipgloss.NewStyle().
			PaddingRight(8).
			PaddingLeft(9).
			Align(lipgloss.Center)

	hoverAddEventStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#7e9cd8")).
				Inherit(addEventStyle)

	hoverCardEventStyle = lipgloss.NewStyle().
				BorderForeground(lipgloss.Color("#7e9cd8")).
				Inherit(cardEventStyle)

	hoverEmptyEventStyle = lipgloss.NewStyle().
				Inherit(emptyEventStyle)

	whiteText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3333"))
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffcc00"))
)

func init() {
	SpinUp()
}

func initialModel(events []Event) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	eventMatrix := CreateEventMatrix(events)
	m := Model{
		spinner:     s,
		events:      events,
		selected:    make(map[Point]struct{}),
		keys:        keys,
		help:        help.New(),
		eventMatrix: eventMatrix,
		mode:        "calendar",
		inputs:      make([]textinput.Model, 8),
		validFields: make([]bool, 8),
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle

		m.inputs[i] = t
		m.validFields[i] = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, textinput.Blink, m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.mode == "calendar" {
		m.keys.Flip.SetEnabled(true)
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.help.Width = msg.Width

		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "?":
				m.help.ShowAll = !m.help.ShowAll

			case "up", "k":
				if m.cursor.y > 0 {
					m.cursor.y--
				}

			case "down", "j":
				if m.cursor.y < EventRowCount(m.events) && m.eventMatrix[m.cursor.y+1][m.cursor.x].Summary != "" {
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
				m.showLocation = !m.showLocation

			case " ", "enter":
				_, ok := m.selected[Point{x: m.cursor.x, y: m.cursor.y}]
				if ok {
					delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
				} else {
					m.mode = "forms"
					m.focusIndex = 0
					event := m.eventMatrix[m.cursor.y][m.cursor.x]
					if event.Summary == "+" {
						m.inputs[Summary].SetValue("")
					} else {
						m.inputs[Summary].SetValue(event.Summary)
					}
					if event.Start.Date == "" {
						m.inputs[Date].SetValue(NewEventDate(m.cursor.x))
					} else {
						m.inputs[Date].SetValue(event.Start.Date)

					}
					m.inputs[StartTime].SetValue(event.Start.DateTime.Format("15:04"))
					m.inputs[EndTime].SetValue(event.End.DateTime.Format("15:04"))
					m.inputs[Location].SetValue(event.Location)
					m.inputs[Id].SetValue(event.Id)
					m.selected[Point{x: m.cursor.x, y: m.cursor.y}] = struct{}{}
				}
			}
		}
	}
	if m.mode == "loading" {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	if m.mode == "forms" {
		m.keys.Flip.SetEnabled(false)
		if m.focusIndex < len(m.inputs)-2 {
			m.keys.Quit.SetEnabled(false)
			m.keys.Help.SetEnabled(false)
		} else {
			m.keys.Quit.SetEnabled(true)
			m.keys.Help.SetEnabled(true)
		}

		var newEvent Event
		newEvent.Summary = "+"
		if m.eventMatrix[m.cursor.y][m.cursor.y] == newEvent {
			m.newEvent = true
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "q":
				if m.focusIndex < len(m.inputs)-2 {
					break
				} else {
					return m, tea.Quit
				}
			case "?":
				if m.focusIndex < len(m.inputs)-2 {
					break
				} else {
					m.help.ShowAll = !m.help.ShowAll
				}
			case "ctrl+c":
				return m, tea.Quit
			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()
				if s == "enter" && m.focusIndex == len(m.inputs)-2 {
					if FormsValidation(m.inputs, &m.validFields) {
						return m, nil
					}
					var err error
					var currentEvent Event
					currentEvent.Id = m.inputs[Id].Value()
					currentEvent.Start.Date = m.inputs[Date].Value()

					formatedStartTime := fmt.Sprintf("%sT%s:00-06:00", m.inputs[Date].Value(), m.inputs[StartTime].Value())
					currentEvent.Start.DateTime, err = time.Parse(time.RFC3339, formatedStartTime)
					if err != nil {
						log.Fatal(err)
					}
					formatedEndTime := fmt.Sprintf("%sT%s:00-06:00", m.inputs[Date].Value(), m.inputs[EndTime].Value())
					currentEvent.End.DateTime, err = time.Parse(time.RFC3339, formatedEndTime)
					if err != nil {
						log.Fatal(err)
					}

					currentEvent.Summary = m.inputs[Summary].Value()
					currentEvent.Location = m.inputs[Location].Value()

					if m.newEvent == true {
						err = PostEvent(currentEvent)
					} else {
						err = UpdateEvent(currentEvent)
					}
					if err != nil {
						log.Println("Failed to update Event:", err)
					}
					m.events = GetEvents()
					m.eventMatrix = CreateEventMatrix(m.events)
					m.newEvent = false
					delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
					m.mode = "calendar"
					return m, nil
				} else if s == "enter" && m.focusIndex == len(m.inputs)-1 {
					for i := range m.validFields {
						m.validFields[i] = true
					}
					m.newEvent = false
					m.areYouSure = false
					delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
					m.mode = "calendar"
				} else if s == "enter" && m.focusIndex == len(m.inputs) {
					if !m.areYouSure {
						m.areYouSure = true
					} else {
						DeleteEvent(m.eventMatrix[m.cursor.y][m.cursor.x])
						m.cursor.y -= 1
						m.events = GetEvents()
						m.eventMatrix = CreateEventMatrix(m.events)
						for i := range m.validFields {
							m.validFields[i] = true
						}
						m.newEvent = false
						m.areYouSure = false
						delete(m.selected, Point{x: m.cursor.x, y: m.cursor.y})
						m.mode = "calendar"

					}
				}

				if s == "up" || s == "shift+tab" {
					m.focusIndex--
				}

				if s == "down" {
					m.focusIndex++
				}

				if m.focusIndex > len(m.inputs)-1 && m.newEvent {
					m.focusIndex = 0
				} else if m.focusIndex < 0 && m.newEvent {
					m.focusIndex = len(m.inputs) - 1
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
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd

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
	switch m.mode {
	case "forms":
		labels := []string{"Event:", "Date:", "Start Time:", "End Time:", "Location:", "Id: "}
		for i := range labels {
			if !m.validFields[i] {
				s += errorStyle.Render(labels[i] + " Invalid field")
				s += fmt.Sprintf("\n%s", m.inputs[i].View())
			} else {
				s += fmt.Sprintf("%s\n%s", labels[i], m.inputs[i].View())
			}
			if i < len(m.inputs)-1 {
				s += "\n"
			}
		}

		submitButton := &blurredSubmitButton
		if m.focusIndex == len(m.inputs)-2 {
			submitButton = &focusedSubmitButton
		}

		cancelButton := &blurredCancelButton
		if m.focusIndex == len(m.inputs)-1 {
			cancelButton = &focusedCancelButton
		}

		var deleteButton *string
		if m.newEvent {
			deleteButton = &grayBlurredDeleteButton
		} else {
			deleteButton = &blurredDeleteButton
			if m.focusIndex == len(m.inputs) {
				deleteButton = &focusedDeleteButton
			}

		}
		var b strings.Builder
		fmt.Fprintf(&b, "\n\n%s %s %s \n\n", *submitButton, *cancelButton, *deleteButton)
		s += b.String()
		if m.areYouSure {
			s += warningStyle.Render("Are you sure?")
		}
	case "loading":
		s += fmt.Sprintf("Loading %s", m.spinner.View())
	case "calendar":
		s += "\n"

		styledDays := GetDaysStartingToday()
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
						if !m.showLocation {
							start := event.Start.DateTime.Format("15:04")
							end := event.End.DateTime.Format("15:04")
							rowEventsTitle = append(rowEventsTitle, hoverCardEventStyle.Render(event.Summary+"\n"+start+"-"+end))
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
						rowEventsTitle = append(rowEventsTitle, cardEventStyle.Render(Truncate(event.Summary, 25, true)))
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
	s += m.help.View(m.keys)
	return s
}
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Flip},
		{k.Help, k.Quit},
	}
}

func main() {
	RefreshOauth()
	events := GetEvents()
	p := tea.NewProgram(initialModel(events))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
