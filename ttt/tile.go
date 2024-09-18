package ttt

type Tile uint

const (
	EmptyTile Tile = iota
	CircleTile
	CrossTile
)

type TileSymbol rune

const (
	EmptySymbol  TileSymbol = ' '
	CircleSymbol TileSymbol = '●'
	CrossSymbol  TileSymbol = 'Χ'
)

func (t Tile) GetSymbol() TileSymbol {
	switch t {
	case CrossTile:
		return CrossSymbol
	case CircleTile:
		return CircleSymbol
	default:
		return EmptySymbol
	}
}
