package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var color1 = lipgloss.NewStyle().
	Padding(1).
	Foreground(lipgloss.Color("#04c26f"))

var color2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#c20404"))

var color3 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("63"))

var color4 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#23469e"))

type model struct {
	items    []NewsItem
	cursor   int
	selected *NewsItem
	viewport viewport.Model
	width    int
	height   int
}

func main() {
	// get source
	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide news source ...")
		return
	}

	// Get the first argument
	source := os.Args[1]

	newsItems, err := get_news(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	m := model{
		items: newsItems,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func renderItems(items []NewsItem, cursor int) string {
	var s string
	welcome := color1.Render("List of News from Source Hackernews")
	s += welcome
	s += "\n\n"
	for i, item := range items {
		cursorStr := " "
		if i == cursor {
			cursorStr = ">>> "
		}
		s += fmt.Sprintf("%s | %s\n", color2.Render(cursorStr), item.Title)
	}
	return s
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport = viewport.New(msg.Width, msg.Height-5)
		m.viewport.SetContent(renderItems(m.items, m.cursor))
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.viewport.SetContent(renderItems(m.items, m.cursor))
			}
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.viewport.SetContent(renderItems(m.items, m.cursor))
			}
		case "enter":
			m.selected = &m.items[m.cursor]
			m.viewport.SetContent(fmt.Sprintf("Title: %s\n\n%s", color3.Render(m.selected.Title), m.selected.Content))
		case "q", "ctrl+c":
			if m.selected != nil {
				m.selected = nil
				m.viewport.SetContent(renderItems(m.items, m.cursor))
			} else {
				return m, tea.Quit
			}
		case "pgup":
			m.viewport.LineUp(5)
		case "pgdown":
			m.viewport.LineDown(5)
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.selected != nil {
		url := fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", m.selected.URL, color4.Render(m.selected.URL))
		return fmt.Sprintf("Title: %s\n\nRead on Web: %s\n\n%s\n\n%s", color3.Render(m.selected.Title), url, m.viewport.View(), color2.Render("Press 'q' to quit."))
	}

	s := m.viewport.View()
	s += "\n"
	s += color1.Render("Press 'q' to quit. Use arrow keys to navigate. Use PgUp/PgDn to scroll.")
	s += "\n"

	return s
}
