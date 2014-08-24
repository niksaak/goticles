package bnticles

import (
	"fmt"
	"github.com/niksaak/goticles/vect"
)

type BadSideIdError Side

func (q BadSideIdError) Error() string {
	return fmt.Sprintf("bad quadrant: %d", q)
}

// Side is a diagonal side of a BB.
type Side int

// helper constants for ease of construction
const (
	north, east Side = 0, 0
	west             = 1
	south            = 2
)

// Valid values of side type
const (
	NorthEast Side = north | east
	NorthWest      = north | west
	SouthEast      = south | east
	SouthWest      = south | west
)

func (s Side) Vect() vect.V {
	switch s {
	case NorthEast:
		return vect.V{1, 1}
	case SouthEast:
		return vect.V{1, -1}
	case SouthWest:
		return vect.V{-1, -1}
	case NorthWest:
		return vect.V{-1, 1}
	default:
		panic(BadSideIdError(s))
	}
}

func (s *Side) String() string {
	var ret string
	if *s & south != 0 {
		ret += "South"
	} else {
		ret += "North"
	}
	if *s & west != 0 {
		ret += "West"
	} else {
		ret += "East"
	}
	return ret
}

// The BB type represents an axis-aligned bounding box, storing its topleft and
// bottomright corners.
type BB [2]vect.V

func (b BB) IsValid() bool {
	return b[0].X < b[1].X && b[0].Y < b[1].Y
}

func (b BB) Center() vect.V {
	return vect.V{
		X: 0.5 * (b[0].X + b[1].X),
		Y: 0.5 * (b[0].Y + b[1].Y),
	}
}

func (b BB) Size() float64 {
	return b[1].X - b[0].X
}

func (b BB) Query(v vect.V) bool {
	return v.X >= b[0].X && v.Y <= b[0].Y &&
		v.X < b[1].X && v.Y > b[1].Y
}

func (b BB) Side(side Side) BB {
	center := b.Center()
	switch side {
	case NorthEast:
		return BB{
			{center.X, b[0].Y},
			{b[1].X, center.Y},
		}
	case SouthEast:
		return BB{
			center,
			b[1],
		}
	case SouthWest:
		return BB{
			{b[0].X, center.Y},
			{center.X, b[1].Y},
		}
	case NorthWest:
		return BB{
			b[0],
			center,
		}
	default:
		panic(BadSideIdError(side))
	}
}
