package bnticles

import (
	//"fmt"
	"github.com/niksaak/goticles"
	"github.com/niksaak/goticles/vect"
)

type Space struct {
	time        float64
	Particles   []goticles.P
	bnParticles [4][]particle
}

func New() *Space {
	return new(Space)
}

func (s *Space) Particle(id int) *goticles.P {
	if id >= len(s.Particles) {
		return nil
	}
	return &s.Particles[id]
}

func (s *Space) MkParticle(mass float64) *goticles.P {
	id := len(s.Particles)
	s.Particles = append(s.Particles, goticles.P{
		Id:          id,
		Mass:        mass,
	})
	for i := range s.bnParticles {
		s.bnParticles[i] = append(s.bnParticles[i], particle{})
	}
	return &s.Particles[id]
}

func (s *Space) RmParticle(p *goticles.P) {
	// TODO
}

func (s *Space) Step(dt float64) {
	s.time += dt
	s.evaluate1()
	s.evaluateK(dt / 2, 1)
	s.evaluateK(dt / 2, 2)
	s.evaluateK(dt,     3)
	s.applyState()
}

func (s *Space) evaluate1() {
	for i, p := range s.Particles {
		s.bnParticles[0][i] = particle{
			position: p.Position,
			velocity: p.Velocity,
			mass:     p.Mass,
		}
	}
}

func (s *Space) evaluateK(dt float64, k int) {
	q := NewQuad(BB{{-1, 1}, {1, -1}})
	copy(s.bnParticles[k], s.bnParticles[0])
	q.insertSlice(s.bnParticles[k])
	for i := range s.bnParticles[k] {
		p := &s.bnParticles[k][i]
		p.position = p.position.Add(s.bnParticles[k-1][i].velocity.Mul(dt))
		p.velocity = p.velocity.Sub(p.treeForce(q))
		p.mass = s.bnParticles[k-1][i].mass
	}
}

func (s *Space) applyState() {
	for i := range s.Particles {
		p := &s.Particles[i]
		p.Position = rk4mean(
			s.bnParticles[0][i].position,
			s.bnParticles[1][i].position,
			s.bnParticles[2][i].position,
			s.bnParticles[3][i].position)
		p.Velocity = rk4mean(
			s.bnParticles[0][i].velocity,
			s.bnParticles[1][i].velocity,
			s.bnParticles[2][i].velocity,
			s.bnParticles[3][i].velocity)
	}
}

func rk4mean(k1, k2, k3, k4 vect.V) vect.V {
	return vect.V{
		X: (1.0/6.0) * (k1.X + 2*(k2.X + k3.X) + k4.X),
		Y: (1.0/6.0) * (k1.Y + 2*(k2.Y + k3.Y) + k4.Y),
	}
}

