package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

type Styles struct {
	focusedStyle            lipgloss.Style
	cursorStyle             lipgloss.Style
	noStyle                 lipgloss.Style
	blurredStyle            lipgloss.Style
	blurredSubmitButton     string
	focusedSubmitButton     string
	blurredCancelButton     string
	focusedCancelButton     string
	blurredDeleteButton     string
	focusedDeleteButton     string
	grayBlurredStyle        lipgloss.Style
	grayBlurredDeleteButton string
	cardEventStyle          lipgloss.Style
	dayStyle                lipgloss.Style
	addEventStyle           lipgloss.Style
	emptyEventStyle         lipgloss.Style
	hoverAddEventStyle      lipgloss.Style
	hoverCardEventStyle     lipgloss.Style
	hoverEmptyEventStyle    lipgloss.Style
	whiteText               lipgloss.Style
	errorStyle              lipgloss.Style
	warningStyle            lipgloss.Style
}
type color struct {
	primary   string
	secondary string
	warning   string
	error     string
}

var colors color

func SetStyles() Styles {
	colors.primary = os.Getenv("COLOR_PRIMARY")
	colors.warning = os.Getenv("COLOR_WARNING")
	colors.error = os.Getenv("COLOR_ERROR")
	w, _, _ := term.GetSize(os.Stdout.Fd())
	myStyles := Styles{}
	myStyles.focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.primary))
	myStyles.cursorStyle = myStyles.focusedStyle
	myStyles.noStyle = lipgloss.NewStyle()

	myStyles.blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.primary))
	myStyles.blurredSubmitButton = fmt.Sprintf("[ %s ]", myStyles.blurredStyle.Render("Submit"))
	myStyles.focusedSubmitButton = myStyles.focusedStyle.Render("[ Submit ]")

	myStyles.blurredCancelButton = fmt.Sprintf("[ %s ]", myStyles.blurredStyle.Render("Cancel"))
	myStyles.focusedCancelButton = myStyles.focusedStyle.Render("[ Cancel ]")

	myStyles.blurredDeleteButton = fmt.Sprintf("[ %s ]", myStyles.blurredStyle.Render("Delete"))
	myStyles.focusedDeleteButton = myStyles.focusedStyle.Render("[ Delete ]")

	myStyles.grayBlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	myStyles.grayBlurredDeleteButton = myStyles.grayBlurredStyle.Render("[ Delete ]")

	myStyles.dayStyle = lipgloss.NewStyle().
		Width((w / 8) + 2).
		Align(lipgloss.Center)

	myStyles.addEventStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, false, true).
		Width(w / 8).
		Height(1).
		Align(lipgloss.Center)
	myStyles.cardEventStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true, true, false, true).
		Width(w / 8).
		Height(5).
		Align(lipgloss.Center)

	myStyles.emptyEventStyle = lipgloss.NewStyle().
		Width((w / 8) + 2).
		Height(5).
		Align(lipgloss.Center)

	myStyles.hoverAddEventStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(colors.primary)).
		Inherit(myStyles.addEventStyle)

	myStyles.hoverCardEventStyle = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(colors.primary)).
		Inherit(myStyles.cardEventStyle)

	myStyles.hoverEmptyEventStyle = lipgloss.NewStyle().
		Inherit(myStyles.emptyEventStyle)

	myStyles.whiteText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))

	myStyles.errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.error))
	myStyles.warningStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.warning))
	return myStyles
}
