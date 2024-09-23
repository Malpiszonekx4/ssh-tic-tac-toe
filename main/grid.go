package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/malpiszonekx4/ssh-tic-tac-toe/ttt"
)

type Point struct {
	X int
	Y int
}

func NewGrid(game *ttt.TicTacToe) *Grid {
	return &Grid{
		game:             game,
		selectedCell:     &Point{},
		selectionVisible: new(bool),
	}
}

type Grid struct {
	game             *ttt.TicTacToe
	selectedCell     *Point
	selectionVisible *bool
}

type ToggleCursorBlinkMsg time.Time

func tickSelectionBlink() tea.Cmd {
	return tea.Every(time.Second/2, func(t time.Time) tea.Msg {
		return ToggleCursorBlinkMsg(t)
	})
}

func (m Grid) Init() tea.Cmd {
	return tickSelectionBlink()
}

func (m Grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d", "right":
			m.selectedCell.X = max(0, min(len(m.game.Map[0])-1, m.selectedCell.X+1))
			*m.selectionVisible = true
		case "a", "left":
			m.selectedCell.X = max(0, min(len(m.game.Map[0])-1, m.selectedCell.X-1))
			*m.selectionVisible = true
		case "w", "up":
			m.selectedCell.Y = max(0, min(len(m.game.Map[0])-1, m.selectedCell.Y-1))
			*m.selectionVisible = true
		case "s", "down":
			m.selectedCell.Y = max(0, min(len(m.game.Map[0])-1, m.selectedCell.Y+1))
			*m.selectionVisible = true
		case "enter", " ":
			*m.selectionVisible = true
			return m, tea.Quit
		}
	case ToggleCursorBlinkMsg:
		*m.selectionVisible = !*m.selectionVisible
		return m, tickSelectionBlink()
	}
	return m, nil
}

func (m Grid) View() string {
	var builder strings.Builder

	rows := len(m.game.Map)
	columns := len(m.game.Map[0])

	// "┌───┬───┬───┐"
	writeHeader(&builder, columns)

	for i := range rows {
		var selectedCell *int = nil

		if m.selectionVisible != nil && *m.selectionVisible && i == int(m.selectedCell.Y) {
			selectedCell = &m.selectedCell.X
		}

		// "│   │   │   │"
		writeRow(&builder, m.game.Map[i], selectedCell)

		if i+1 == rows {
			break
		}

		// "├───┼───┼───┤"
		writeSeparator(&builder, columns)

	}
	// "└───┴───┴───┘"
	writeFooter(&builder, columns)

	return builder.String()
}

func writeHeader(builder *strings.Builder, columns int) {
	builder.WriteRune('┌')
	for j := 0; j < columns; j++ {
		builder.WriteString("───")
		if j < columns-1 {
			builder.WriteRune('┬')
		}
	}
	builder.WriteRune('┐')
	builder.WriteRune('\n')
}

func writeRow(builder *strings.Builder, row []ttt.Tile, selectedColumn *int) {
	for i, tile := range row {
		builder.WriteRune('│')
		builder.WriteRune(' ')
		if selectedColumn != nil && i == *selectedColumn {
			builder.WriteString("\u001b[37;47m")
			builder.WriteRune(rune(tile.GetSymbol()))
			builder.WriteString("\u001b[0m")
		} else {
			builder.WriteRune(rune(tile.GetSymbol()))
		}
		builder.WriteRune(' ')
	}
	builder.WriteString("│\n")
}

func writeSeparator(builder *strings.Builder, columns int) {
	builder.WriteRune('├')
	for j := 0; j < columns; j++ {
		builder.WriteString("───")
		if j < columns-1 {
			builder.WriteRune('┼')
		}
	}
	builder.WriteRune('┤')
	builder.WriteRune('\n')
}

func writeFooter(builder *strings.Builder, columns int) {
	builder.WriteRune('└')
	for j := 0; j < columns; j++ {
		builder.WriteString("───")
		if j < columns-1 {
			builder.WriteRune('┴')
		}
	}
	builder.WriteRune('┘')
	builder.WriteRune('\n')
}
