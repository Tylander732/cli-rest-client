package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylander732/cli-rest-client/pkg/http"
)

// Text area size settings
const (
	//Starting number of text areas
	initialInputs = 2
	helpHeight    = 5

	uriMinHeight      = 1
	uriMinWidth       = 1
	treeMinHeight     = 1
	treeMinWidth      = 1
	dataMinHeight     = 1
	dataMinWidth      = 1
	responseMinHeight = 1
	responseMinWidth  = 1
)

var (
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	cursorLineStyle  = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	endOfBufferStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("235"))
	focusedPlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	focusedBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238"))
	blurredBorderStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
)

type keymap = struct {
	//Current available keybinding types
	next, prev, quit, send key.Binding
}

// Settings for the text areas?
type model struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

func newModel() model {
	m := model{
		inputs: make([]textarea.Model, initialInputs),
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
			//TODO: Set this as enter instead, but for some reason ctrl+enter doesn't work
			send: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "send"),
			),
		},
	}
	for i := 0; i < initialInputs; i++ {
		m.inputs[i] = newTextarea()
	}
	m.inputs[m.focus].Focus()
	return m
}

// (m model) is the receiver of the method. It means that the 'Init'
// method is accosiated with a type 'model'
// I'm attaching this func to a struct, kinda like a method in a class?
// When I call init on a model struct, make the textarea blink (this is a tea cmd)
func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// Msg is input from the user
// Update handles reading input from the user as well as performing
// UI updates
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// If the key is the Quit key, return the model and tea.Quit
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, tea.Quit

		// Next keymap
		case key.Matches(msg, m.keymap.next):
			// Prior to moving to the next text area, blur the current one
			m.inputs[m.focus].Blur()

			// Move focus to the next text area
			m.focus++

			// If we've reached the end of the available areas, go back to first text area
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)

		// Prev keymap
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)

		case key.Matches(msg, m.keymap.send):
			log.Println(m.inputs[m.focus].Value())

			//TODO: These values need to be updated
			response := http.GetRequest("", "http://localhost:8080/items")
			m.inputs[1].SetValue(response)
		}

	// This will always happen once initially, it then happens whenever the user resizes the terminal window
	case tea.WindowSizeMsg:
		//TODO: possible initial height and width bug here
		m.height = msg.Height
		m.width = msg.Width
	}

	for i := range m.inputs {
		// inputs.Update() is the bubletea update loop
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// This is used to change the size of all the text areas as they're created
// Maybe it will be useful for something else
func (m *model) sizeInputs() {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

// Attaches View method to the model struct?
// What is View actually doing?
func (m model) View() string {

	// Setting up a help object that contains all help messages for my keybindings
	help := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		m.keymap.quit,
		m.keymap.send,
	})

	// Setup a string slice that contains all views? Why would a string slice contain the views?
	var views []string
	for i := range m.inputs {
		// append each of the inputs on the model to the views slice
		// inputs.View() renders the text area
		// View() Returns a string, so all UI elements must just be strings?
		views = append(views, m.inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)

	log.Println("Hello")

	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(0)
	}
}
