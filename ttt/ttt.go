package ttt

type Row uint

const (
	Row1 Row = iota
	Row2
	Row3
)

type Column uint

const (
	Column1 Column = iota
	Colmun2
	Colmun3
)

type TicTacToe struct {
	Map           [][]Tile
	CurrentPlayer Player
}

func CreateGame() *TicTacToe {
	game := TicTacToe{
		CurrentPlayer: getRandomPlayer(),
		Map: [][]Tile{
			{EmptyTile, EmptyTile, EmptyTile},
			{EmptyTile, EmptyTile, EmptyTile},
			{EmptyTile, EmptyTile, EmptyTile},
		},
	}

	return &game
}

type MoveResult int

const (
	Ok MoveResult = iota
	Win
	NoMoreMoves
	InvalidMove
)

func (t *TicTacToe) Move(row Row, column Column) MoveResult {
	if t.Map[row][column] != EmptyTile {
		return InvalidMove
	}

	t.Map[row][column] = t.CurrentPlayer.getTileState()

	if noMoreMoves(t.Map) {
		return NoMoreMoves
	}

	if checkWinInRow(t.Map) || checkWinInColumn(t.Map) || checkWinInDiagonal(t.Map) {
		return Win
	}

	t.nextPlayer()

	return Ok
}

func (t *TicTacToe) nextPlayer() {
	switch t.CurrentPlayer {
	case Circle:
		t.CurrentPlayer = Cross
	case Cross:
		t.CurrentPlayer = Circle
	}
}

func noMoreMoves(gameMap [][]Tile) bool {
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

func checkWinInColumn(gameMap [][]Tile) bool {
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

func checkWinInRow(gameMap [][]Tile) bool {
	for _, row := range gameMap {
		if row[0] != EmptyTile && row[0] == row[1] && row[1] == row[2] {
			return true
		}
	}

	return false
}

func checkWinInDiagonal(gameMap [][]Tile) bool {
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
