package goticles

import (
	"testing"
)

func mkParticle() Particle {
	p := Particle{ id: 0 }
	p.SetMass(1)
	return p
}

func mkSpace() Space {
	return Space{
		particles: []Particle{
			mkParticle(),
		},
		force: func(s State, t float64) Vect {
			return Vect{0, 0}
		},
	}
}

func TestSpace(t *testing.T) {
	space := mkSpace()
	particle := &space.particles[0]
	particle.ApplyImpulse(Vect{1, 0})
	step := 1.0/60
	for i := 0.0; i < 1.01; i += step {
		t.Logf("t: %.4f, p: %v", i, particle)
		space.Integrate(step)
	}
}

