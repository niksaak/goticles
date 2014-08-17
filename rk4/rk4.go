package goticles

import (
	"github.com/niksaak/goticles/vect"
	"math"
)

type SpaceRK4 struct {
	Time         float64
	Particles    []Particle
	positions    [][4]vect.V
	velocities   [][4]vect.V
	masses       []float64
}

func MkSpaceRK4() *SpaceRK4 {
	return new(SpaceRK4)
}

func (s *SpaceRK4) Particle(id int) *Particle {
	return &s.Particles[id]
}

func (s *SpaceRK4) MkParticle(mass float64) *Particle {
	id := len(s.Particles)
	s.Particles = append(s.Particles, Particle{
		Id:          id,
		Mass:        mass,
		massInverse: 1.0 / mass,
	})
	s.positions = append(s.positions, [4]vect.V{})
	s.velocities = append(s.velocities, [4]vect.V{})
	return &s.Particles[id]
}

func (s *SpaceRK4) Step(dt float64) {
	s.evaluate1()
	s.evaluateK(dt/2, 1)
	s.evaluateK(dt/2, 2)
	s.evaluateK(dt, 3)
	s.applyState(dt)
	s.Time += dt
}

func (s *SpaceRK4) evaluate1() {
	particleCount := len(s.Particles)

	// check integration arrays size
	if len(s.positions) != particleCount {
		s.positions = make([][4]vect.V, particleCount)
	}
	if len(s.velocities) != particleCount {
		s.velocities = make([][4]vect.V, particleCount)
	}
	if len(s.masses) != particleCount {
		s.masses = make([]float64, particleCount)
	}

	// get state
	for i, p := range s.Particles {
		s.positions[i][0] = p.Position
		s.velocities[i][0] = p.Velocity
		s.masses[i] = p.Mass
	}
}

func (s *SpaceRK4) evaluateK(dt float64, k int) {
	for i := range s.positions {
		position := &s.positions[i][k]
		position.X = s.positions[i][0].X + s.velocities[i][k-1].X*dt
		position.Y = s.positions[i][0].Y + s.velocities[i][k-1].Y*dt

		s.velocities[i][k] = s.velocities[i][0]
	}
	for i := range s.positions {
		position := &s.positions[i][k]
		m1 := s.masses[i]
		for j := i + 1; j < len(s.positions); j++ {
			m2 := s.masses[j]

			dx := position.X - s.positions[j][k].X
			dy := position.Y - s.positions[j][k].Y
			if dx == 0 && dy == 0 {
				continue
			}

			dSquared := dx*dx + dy*dy
			distance := math.Sqrt(dSquared)
			mag := dSquared + distance*0

			fX := G * m1 * m2 * dx / mag
			fY := G * m1 * m2 * dy / mag

			s.velocities[i][k].X -= fX * dt
			s.velocities[i][k].Y -= fY * dt

			s.velocities[j][k].X += fX * dt
			s.velocities[j][k].Y += fY * dt
		}
	}
}

func (s *SpaceRK4) applyState(dt float64) {
	for i := range s.Particles {
		p := &s.Particles[i]
		p.Position = rk4Mean(dt, s.positions[i])
		p.Velocity = rk4Mean(dt, s.velocities[i])
	}
}

func rk4Mean(dt float64, vec [4]vect.V) vect.V {
	const FRAC = 1.0 / 6.0
	return vect.V{
		X: FRAC * (vec[0].X + 2*(vec[1].X+vec[2].X) + vec[3].X),
		Y: FRAC * (vec[0].Y + 2*(vec[1].Y+vec[2].Y) + vec[3].Y),
	}
}
