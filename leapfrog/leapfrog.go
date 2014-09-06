package leapfrog

import (
	"github.com/niksaak/goticles"
	"github.com/niksaak/goticles/vect"
)

const (
	G           = 6.67384e-11
	TRESHOLD    = 2e-3
	TRESHOLD_SQ = TRESHOLD * TRESHOLD
)

type Space struct {
	time      float64
	Particles []goticles.P
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
		Id:   id,
		Mass: mass,
	})
	return &s.Particles[id]
}

func (s *Space) RmParticle(id int) {
	// TODO
}

func force(p, q *goticles.P) vect.V {
	distV := p.Position.Sub(q.Position)
	distSq := distV.LenSq()
	if distSq < TRESHOLD_SQ {
		return vect.V{}
	}
	distU := distV.Ulen()
	return distU.Mul(G * p.Mass * q.Mass / distSq).Neg()
}

func particleLeapfrog(p *goticles.P, force vect.V, dt float64) {
	p.Position = p.Position.Add(p.Velocity.Mul(dt)).Add(
		p.Acceleration.Mul(dt * dt).Mul(0.5))
	p.Velocity = p.Velocity.Add(p.Acceleration.Add(force).Mul(0.5).Mul(dt))
	p.Acceleration = force
}

func (s *Space) Step(dt float64) {
	ln := len(s.Particles)
	for i := range s.Particles {
		p := &s.Particles[i]
		for j := i + 1; j < ln; j++ {
			q := &s.Particles[j]
			pqForce := force(p, q)
			particleLeapfrog(p, pqForce, dt)
			particleLeapfrog(q, pqForce.Neg(), dt)
		}
	}
}
