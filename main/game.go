package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/malpiszonekx4/ssh-tic-tac-toe/ttt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/samber/lo"
)

type PlayerSymbol rune

const (
	EmptySymbol  PlayerSymbol = ' '
	CircleSymbol PlayerSymbol = '●'
	CrossSymbol  PlayerSymbol = 'Χ'
)

func getSymbolForTile(t ttt.Tile) PlayerSymbol {
	switch t {
	case ttt.CrossTile:
		return CrossSymbol
	case ttt.CircleTile:
		return CircleSymbol
	default:
		return EmptySymbol
	}
}

type Model struct {
	Game                    *ttt.TicTacToe
	moveInput               *textinput.Model
	playAgainConfirmation   *confirmation.Model
	displayPlayAgainConfirm bool
	statusMsg               string
}

func initialModel() Model {
	model := Model{
		// moveInput:     createInput(),
		Game:                    ttt.CreateGame(),
		playAgainConfirmation:   confirmation.NewModel(confirmation.New("Wanna play again?", confirmation.Yes)),
		displayPlayAgainConfirm: false,
		statusMsg:               "",
	}
	model.moveInput = createInput(&model)
	return model
}

func createInput(model *Model) *textinput.Model {
	moveInput := textinput.New("Your move")
	moveInput.Validate = func(input string) error {
		matched, _ := regexp.MatchString("^(?:[abcABC][123])|(?:[123][abcABC])$", input)
		if !matched {
			return fmt.Errorf("invalid move")
		}

		row, column := parsePlayerMove(input)
		if model.Game.Map[row][column] != ttt.EmptyTile {
			return fmt.Errorf("tile not empty")
		}

		return nil
	}
	moveInput.Placeholder = "e.g. A2"

	return textinput.NewModel(moveInput)
}

func stringifyTileState(ts [][]ttt.Tile) [][]string {
	return lo.Map(ts, func(x []ttt.Tile, _ int) []string {
		return lo.Map(x, func(y ttt.Tile, _ int) string {
			return string(getSymbolForTile(y))
		})
	})
}

func (m Model) Init() tea.Cmd {
	m.moveInput.Init()
	m.playAgainConfirmation.Init()

	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			fmt.Println("Bye!")
			return m, tea.Quit
		default:
			if m.displayPlayAgainConfirm {
				_, cmd := m.playAgainConfirmation.Update(msg)
				if cmd != nil && cmd() == tea.Quit() {
					playAgain, err := m.playAgainConfirmation.Value()
					if err != nil {
						fmt.Println("cannot get playAgainConfirmation.Value()")
					}
					if playAgain {
						m = initialModel()
						m.Init()
					} else {
						return m, tea.Quit
					}
				}
			} else {
				_, cmd := m.moveInput.Update(msg)
				if cmd != nil && cmd() == tea.Quit() {
					val, err := m.moveInput.Value()
					if err != nil {
						fmt.Println("Error getting input value")
					}

					row, column := parsePlayerMove(val)

					moveResult := m.Game.Move(row, column)

					switch moveResult {
					case ttt.Win:
						m.displayPlayAgainConfirm = true
						m.statusMsg = "The winner is " + m.Game.CurrentPlayer.GetShapeName()
					case ttt.NoMoreMoves:
						m.displayPlayAgainConfirm = true
						m.statusMsg = "It's a tie"
					case ttt.Ok:
						m.moveInput = createInput(&m)
						m.moveInput.Init()
					}
				}
			}
		}
	}
	return m, nil
}

func parsePlayerMove(input string) (row ttt.Row, column ttt.Column) {
	char1 := rune(input[0])
	char2 := rune(input[1])

	switch char1 {
	case 'a', 'A':
		column = ttt.Column1
	case 'b', 'B':
		column = ttt.Colmun2
	case 'c', 'C':
		column = ttt.Colmun3
	default:
		tmpRow, err := strconv.Atoi(string(char1))
		if err != nil {
			fmt.Println("Error while parsing input")
		}
		row = ttt.Row(uint32(tmpRow) - 1)
	}

	switch char2 {
	case 'a', 'A':
		column = ttt.Column1
	case 'b', 'B':
		column = ttt.Colmun2
	case 'c', 'C':
		column = ttt.Colmun3
	default:
		tmpRow, err := strconv.Atoi(string(char2))
		if err != nil {
			fmt.Println("Error while parsing input")
		}
		row = ttt.Row(uint32(tmpRow) - 1)
	}

	return
}

func (m Model) View() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Rows(stringifyTileState(m.Game.Map)...).
		BorderRow(true).
		BorderColumn(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	cols := style.Render(strings.Join([]string{"A", "B", "C  "}, "   "))
	rows := style.Render(strings.Join([]string{"1", "2", "3"}, " \n\n"))

	s := lipgloss.JoinVertical(lipgloss.Right, lipgloss.JoinHorizontal(lipgloss.Center, rows, t.Render()), cols) + "\n"

	if !m.displayPlayAgainConfirm {
		s += m.moveInput.View()
	} else {
		s += m.statusMsg + "\n"
		s += m.playAgainConfirmation.View()
	}

	return s + "\n"
}

// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }
