package main

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/samber/lo"
)

type TileState rune

const (
	EmptyTile  TileState = ' '
	CircleTile TileState = '●'
	CrossTile  TileState = 'Χ'
)

type Player bool

const (
	X Player = false
	O Player = true
)

func (p Player) getTile() TileState {
	if p == X {
		return CrossTile
	} else {
		return CircleTile
	}
}

func (p Player) getName() string {
	switch p {
	case X:
		return "cross"
	case O:
		return "circle"
	default:
		return ""
	}
}

type Model struct {
	Map                     [][]TileState
	moveInput               *textinput.Model
	currentPlayer           Player
	playAgainConfirmation   *confirmation.Model
	displayPlayAgainConfirm bool
	statusMsg               string
}

func initialModel() Model {
	model := Model{
		Map: [][]TileState{
			{EmptyTile, EmptyTile, EmptyTile},
			{EmptyTile, EmptyTile, EmptyTile},
			{EmptyTile, EmptyTile, EmptyTile},
		},
		// moveInput:     createInput(),
		currentPlayer:           getRandomPlayer(),
		playAgainConfirmation:   confirmation.NewModel(confirmation.New("Wanna play again?", confirmation.Yes)),
		displayPlayAgainConfirm: false,
		statusMsg:               "",
	}
	model.moveInput = createInput(&model)
	return model
}

func getRandomPlayer() Player {
	if rand.N(2) == 0 {
		return X
	} else {
		return O
	}
}

func (m *Model) nextPlayer() {
	if m.currentPlayer == X {
		m.currentPlayer = O
	} else {
		m.currentPlayer = X
	}
}

func createInput(model *Model) *textinput.Model {
	moveInput := textinput.New("Your move")
	moveInput.Validate = func(input string) error {
		matched, _ := regexp.MatchString("^(?:[abcABC][123])|(?:[123][abcABC])$", input)
		if !matched {
			return fmt.Errorf("invalid move")
		}

		row, column := parsePlayerMove(input)
		if model.Map[row][column] != EmptyTile {
			return fmt.Errorf("tile not empty")
		}

		return nil
	}
	moveInput.Placeholder = "e.g. A2"

	return textinput.NewModel(moveInput)
}

func stringifyTileState(ts [][]TileState) [][]string {
	return lo.Map(ts, func(x []TileState, _ int) []string {
		return lo.Map(x, func(y TileState, _ int) string {
			return string(y)
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

					m.Map[row][column] = m.currentPlayer.getTile()

					if checkWinInColumn(m.Map) || checkWinInRow(m.Map) || checkWinInDiagonal(m.Map) {
						m.displayPlayAgainConfirm = true
						m.statusMsg = "The winner is " + m.currentPlayer.getName()
					} else if noMoreMoves(m.Map) {
						m.displayPlayAgainConfirm = true
						m.statusMsg = "It's a tie"
					}

					m.nextPlayer()

					m.moveInput = createInput(&m)
					m.moveInput.Init()
				}
			}
		}
	}
	return m, nil
}

func noMoreMoves(gameMap [][]TileState) bool {
	emptyTiles := 0
	for _, row := range gameMap {
		for _, cell := range row {
			if cell == EmptyTile {
				emptyTiles++
			}
		}
	}
	return emptyTiles == 0
}

func checkWinInColumn(gameMap [][]TileState) bool {
	row1 := gameMap[0]
	row2 := gameMap[1]
	row3 := gameMap[2]

	for _, col := range []int{0, 1, 2} {
		if row1[col] != EmptyTile && row1[col] == row2[col] && row2[col] == row3[col] {
			return true
		}
	}

	return false
}

func checkWinInRow(gameMap [][]TileState) bool {
	for _, row := range gameMap {
		if row[0] != EmptyTile && row[0] == row[1] && row[1] == row[2] {
			return true
		}
	}

	return false
}

func checkWinInDiagonal(gameMap [][]TileState) bool {
	row1 := gameMap[0]
	row2 := gameMap[1]
	row3 := gameMap[2]

	if row1[2] != EmptyTile && row1[2] == row2[1] && row2[1] == row3[0] {
		return true
	}
	if row1[0] != EmptyTile && row1[0] == row2[1] && row2[1] == row3[2] {
		return true
	}

	return false
}

func parsePlayerMove(input string) (row uint32, column uint32) {
	char1 := rune(input[0])
	char2 := rune(input[1])

	switch char1 {
	case 'a', 'A':
		column = 0
	case 'b', 'B':
		column = 1
	case 'c', 'C':
		column = 2
	default:
		tmpRow, err := strconv.Atoi(string(char1))
		if err != nil {
			fmt.Println("Error while parsing input")
		}
		row = uint32(tmpRow) - 1
	}

	switch char2 {
	case 'a', 'A':
		column = 0
	case 'b', 'B':
		column = 1
	case 'c', 'C':
		column = 2
	default:
		tmpRow, err := strconv.Atoi(string(char2))
		if err != nil {
			fmt.Println("Error while parsing input")
		}
		row = uint32(tmpRow) - 1
	}

	return
}

func (m Model) View() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Rows(stringifyTileState(m.Map)...).
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
