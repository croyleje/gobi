package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The default colors and border style is set here via the lipgloss.Color()
// function please https://github.com/charmbracelet/lipgloss for explanation of
// how the colors are set 4,8,24 bit colors are supported.
var (
	termWidth, termHight int
	heading              = lipgloss.NewStyle().Bold(true).Margin(1, 2)
	notSelected          = lipgloss.NewStyle().Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("12")).Width(14).Padding(1).Margin(1)
	selected             = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")).Width(14).Padding(1).Margin(1)
)

// Default commands.
const (
	Suspend   = iota
	Lock      = iota
	Logout    = iota
	Shutdown  = iota
	Reboot    = iota
	Hibernate = iota
	Exit      = iota
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

// Commands passed to your terminal self-explanatory note flags and options are
// passed with in comma separated list inside quotes.
func parse(choice string) {
	switch choice {
	case "Suspend":
		cmd := exec.Command("systemctl", "suspend")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Lock":
		cmd := exec.Command("hyprlock")
		werr := cmd.Run()
		if werr != nil {
			log.Fatal(werr)
		}

	case "Logout":
		cmd := exec.Command("loginctl", "terminate-user", "$USER")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}

	case "Exit":
		cmd := exec.Command("hyprctl", "dispatch", "exit")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Shutdown":
		cmd := exec.Command("systemctl", "poweroff")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Reboot":
		cmd := exec.Command("systemctl", "reboot")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	case "Hibernate":
		cmd := exec.Command("systemctl", "hibernate")
		err := cmd.Run()

		if err != nil {
			log.Fatal(err)
		}
	}
}

// IF your want to change the order the commands are displayed or remove a
// command do it here.
func initialModel() model {
	return model{
		choices:  []string{"Suspend", "Lock", "Logout", "Shutdown", "Reboot", "Hibernate", "Exit"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(nil, tea.EnterAltScreen)
}

// Default key bindings.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		termWidth, termHight = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k", "h", "shift+tab":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "down", "l", "j", "tab":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case strconv.Itoa(Suspend + 1):
			parse("Suspend")
			return m, tea.Quit
		case strconv.Itoa(Lock + 1):
			parse("Lock")
			return m, tea.Quit
		case strconv.Itoa(Logout + 1):
			parse("Logout")
			return m, tea.Quit
		case strconv.Itoa(Shutdown + 1):
			parse("Shutdown")
			return m, tea.Quit
		case strconv.Itoa(Reboot + 1):
			parse("Reboot")
			return m, tea.Quit
		case strconv.Itoa(Hibernate + 1):
			parse("Hibernate")
			return m, tea.Quit
		case "enter":
			parse(m.choices[m.cursor])
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Header text.
	header := heading.Render("gobi v0.1")
	choices := lipgloss.NewStyle().Render()

	for i, choice := range m.choices {
		cursor := " "
		var entry string
		if m.cursor == i {
			cursor = ">"
			entry = selected.Render(fmt.Sprintf("%s %s", cursor, choice))
		} else {
			entry = notSelected.Render(fmt.Sprintf("%s %s", cursor, choice))
		}

		choices = lipgloss.JoinHorizontal(lipgloss.Center, choices, entry)
	}
	choices = lipgloss.JoinVertical(lipgloss.Center, header, choices)

	// Colors and font settings for footer.
	footer := lipgloss.NewStyle().Faint(true).Italic(true)
	choices = lipgloss.JoinVertical(lipgloss.Center, choices, footer.Render("Vim keys or Tab/Shift Tab to select."))
	choices = lipgloss.JoinVertical(lipgloss.Center, choices, footer.Render("Press esc/q to quit."))

	textWidth, textHeight := lipgloss.Size(choices)
	marginW, marginH := (termWidth-textWidth)/2, (termHight-textHeight)/2

	return lipgloss.NewStyle().Margin(marginH, marginW).Render(choices)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
