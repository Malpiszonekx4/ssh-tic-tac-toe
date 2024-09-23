package main

import (
	"fmt"
	"strings"

	"github.com/malpiszonekx4/ssh-tic-tac-toe/ttt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
)

type Model struct {
	Game                    *ttt.TicTacToe
	grid                    *Grid
	playAgainConfirmation   *confirmation.Model
	displayPlayAgainConfirm bool
	statusMsg               string
}

func initialModel() Model {
	game := ttt.CreateGame()
	model := Model{
		grid:                    NewGrid(game),
		Game:                    game,
		playAgainConfirmation:   confirmation.NewModel(confirmation.New("Wanna play again?", confirmation.Yes)),
		displayPlayAgainConfirm: false,
		statusMsg:               getPlayerTurnStatusMsg(game.CurrentPlayer),
	}
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.grid.Init(),
		m.playAgainConfirmation.Init(),
	)
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
				_, cmd := m.grid.Update(msg)
				if cmd != nil && cmd() == tea.Quit() {
					moveResult := m.Game.Move(ttt.Row(m.grid.selectedCell.Y), ttt.Column(m.grid.selectedCell.X))
					switch moveResult {
					case ttt.Win:
						m.displayPlayAgainConfirm = true
						m.statusMsg = m.Game.CurrentPlayer.GetShapeName() + " won!"
					case ttt.NoMoreMoves:
						m.displayPlayAgainConfirm = true
						m.statusMsg = "It's a tie"
					case ttt.Ok:
						m.statusMsg = getPlayerTurnStatusMsg(m.Game.CurrentPlayer)
					}
				}
			}
		}
	case ToggleCursorBlinkMsg:
		_, cmd := m.grid.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	s := lipgloss.JoinVertical(lipgloss.Center, m.grid.View(), m.statusMsg)

	if m.displayPlayAgainConfirm {
		s = lipgloss.JoinVertical(lipgloss.Left, s, m.playAgainConfirmation.View())
	}

	return s + "\n"
}

func getPlayerTurnStatusMsg(player ttt.Player) string {
	builder := strings.Builder{}
	shape := player.GetShapeName()
	builder.WriteString(shape)
	builder.WriteRune('\'')
	if !strings.HasSuffix(shape, "s") {
		builder.WriteRune('s')
	}
	builder.WriteString(" turn")
	return builder.String()
}
