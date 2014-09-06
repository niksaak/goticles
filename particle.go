package goticles

import (
	"github.com/niksaak/goticles/vect"
)

// The P type describes a Particle to be used in particle simulations.
type P struct {
	Id       int
	Position vect.V
	Velocity vect.V
	Acceleration vect.V
	Force    vect.V
	Mass     float64
}
