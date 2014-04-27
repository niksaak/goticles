package goticles

import (
	"fmt"
	"math"
)

type Particle struct {
	id int
	State
}

func (p *Particle) String() string {
	return fmt.Sprintf(
		"#%d @(% 4.4f, % 4.4f) V(% 4.4f, % 4.4f) ->(% 4.4f, % 4.4f) "+
		" M = % 2f(% 2f)",
		p.id,
		p.position.X,
		p.position.Y,
		p.velocity.X,
		p.velocity.Y,
		p.momentum.X,
		p.momentum.Y,
		p.mass,
		p.massInverse,
	)
}

type Space struct {
	time float64
	particles []Particle
}

func MkSpace() *Space {
	space := new(Space)
	return space
}

func (s *Space) Particle(id int) *Particle {
	return &s.particles[id]
}

func (s *Space) AddParticle(mass float64) *Particle {
	s.particles = append(s.particles, Particle{})
	particle := &s.particles[len(s.particles)-1]
	particle.id = len(s.particles)
	particle.SetMass(mass)
	particle.recalculate()
	return particle
}

func (s *Space) ParticleCount() int {
	return len(s.particles)
}

func (s *Space) applyForce(dt float64) {
	for i := range s.particles {
		pt1 := &s.particles[i]
		for j := i + 1; j < len(s.particles); j++ {
			pt2 := &s.particles[j]

			dx := pt1.position.X - pt2.position.X
			dy := pt1.position.Y - pt2.position.Y

			dSquared := dx*dx + dy*dy
			distance := math.Sqrt(dSquared)
			mag := dSquared * distance

			pt1.force.X -= dx * pt2.mass * mag
			pt1.force.Y -= dy * pt2.mass * mag

			pt2.force.X += dx * pt1.mass * mag
			pt2.force.Y += dy * pt1.mass * mag
		}
	}
}

func (s *Space) Integrate(dt float64) {
	s.applyForce(dt)
	for i := range s.particles {
		integrate(&s.particles[i].State, s.time, dt)
	}
	s.time += dt
}

