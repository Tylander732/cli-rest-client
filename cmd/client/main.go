package main

import (
	// "fmt"
	// "os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
	// tea "github.com/charmbracelet/bubbletea"
)

const (
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
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type keymap = struct {
	next, prev, quit key.Binding
}

type model struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
}

func main() {

}
