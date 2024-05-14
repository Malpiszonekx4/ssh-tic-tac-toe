package ttt

import "math/rand/v2"

type Player uint

const (
	Circle Player = iota
	Cross
)

func (p Player) getTileState() Tile {
	switch p {
	case Circle:
		return CircleTile
	case Cross:
		return CrossTile
	default:
		return EmptyTile
	}
}

func (p Player) GetShapeName() string {
	switch p {
	case Cross:
		return "cross"
	case Circle:
		return "circle"
	default:
		return ""
	}
}

func getRandomPlayer() Player {
	if rand.N(2) == 0 {
		return Cross
	} else {
		return Circle
	}
}
