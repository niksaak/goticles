package goticles

import (
	"github.com/niksaak/goticles/vect"
	"fmt"
	"math"
)

var SuperVerbose bool

type Particle struct {
	Id          int
	Position    vect.V
	Velocity    vect.V
	Force       vect.V
	Mass        float64
	massInverse float64
}

func (p *Particle) SetPosition(pos vect.V) {
	p.Position = pos
}

func (p *Particle) SetMass(mass float64) {
	p.Mass = mass
	p.massInverse = 1 / mass
}

func (p *Particle) SetVelocity(vel vect.V) {
	p.Velocity = vel
}

func (p *Particle) ApplyImpulse(imp vect.V) {
	p.SetVelocity(imp.Mul(p.massInverse))
}

func (p *Particle) String() string {
	return fmt.Sprintf(
		"#%4d; X(%2.4f, %2.4f); V(%2.4f, %2.4f)",
		p.Id,
		p.Position.X, p.Position.Y,
		p.Velocity.X, p.Velocity.Y)
}

type Space struct {
	time float64
	Particles []Particle
	positions [][4]vect.V
	velocities [][4]vect.V
	masses []float64
	massInverses []float64
}

func MkSpace() *Space {
	return new(Space)
}

func (s *Space) Particle(id int) *Particle {
	return &s.Particles[id]
}

func (s *Space) MkParticle(mass float64) *Particle {
	id := len(s.Particles)
	s.Particles = append(s.Particles, Particle{
		Id: id,
		Mass: mass,
		massInverse: 1/mass,
	})
	s.positions = append(s.positions, [4]vect.V{})
	s.velocities = append(s.velocities, [4]vect.V{})
	return &s.Particles[id]
}

const G = 6.67384e-11
var _ = math.MaxFloat64 // TODO: remove this

func (s *Space) Step(dt float64) {
	s.evaluate1()
	s.evaluateK(dt/2, 1)
	s.evaluateK(dt/2, 2)
	s.evaluateK(dt, 3)
	s.applyState(dt)
}

func (s *Space) evaluate1() {
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
	if len(s.massInverses) != particleCount {
		s.massInverses = make([]float64, particleCount)
	}

	// get state
	for i, p := range s.Particles {
		s.positions[i][0] = p.Position
		s.velocities[i][0] = p.Velocity
		s.masses[i] = p.Mass
		s.massInverses[i] = p.massInverse
	}
}

func (s *Space) evaluateK(dt float64, k int) {
	for i := range s.positions {
		position := &s.positions[i][k]
		position.X = s.positions[i][0].X + s.velocities[i][k-1].X * dt
		position.Y = s.positions[i][0].Y + s.velocities[i][k-1].Y * dt

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

func (s *Space) applyState(dt float64) {
	for i := range s.Particles {
		p := &s.Particles[i]
		p.Position = rk4Mean(dt, s.positions[i])
		p.Velocity = rk4Mean(dt, s.velocities[i])
	}
}

func rk4Mean(dt float64, vec [4]vect.V) vect.V {
	const FRAC = 1.0 / 6.0
	return vect.V{
		X: /*dt */ FRAC * (vec[0].X + 2*(vec[1].X + vec[2].X) + vec[3].X),
		Y: /*dt */ FRAC * (vec[0].Y + 2*(vec[1].Y + vec[2].Y) + vec[3].Y),
	}
}

