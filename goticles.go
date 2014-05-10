package goticles

import (
	"github.com/niksaak/goticles/vect"
	"fmt"
	"math"
)

type Particle struct {
	Id          int
	Position    vect.V
	dPosition   [4]vect.V
	Momentum    vect.V
	dMomentum   [4]vect.V
	velocity    vect.V
	dvelocity   [4]vect.V
	Force       vect.V
	dForce      [4]vect.V
	Mass        float64
	massInverse float64
}

func (p *Particle) SetPosition(pos vect.V) {
	p.Position = pos
	p.Recalculate()
}

func (p *Particle) SetMomentum(momentum vect.V) {
	p.Momentum = momentum
	p.Recalculate()
}

func (p *Particle) SetForce(force vect.V) {
	p.Force = force
}

func (p *Particle) SetMass(mass float64) {
	p.Mass = mass
	p.massInverse = 1 / mass
}

func (p Particle) Velocity() vect.V {
	return p.velocity
}

func (p *Particle) String() string {
	return fmt.Sprintf(
		"#%4d; X(%2.4f, %2.4f); P(%2.4f, %2.4f); V(%2.4f, %2.4f)",
		p.Id,
		p.Position.X, p.Position.Y,
		p.Momentum.X, p.Momentum.Y,
		p.velocity.X, p.velocity.Y)
}

func (p *Particle) Recalculate() {
	p.velocity.X = p.Momentum.X * p.massInverse
	p.velocity.Y = p.Momentum.Y * p.massInverse
}

func (p *Particle) dRecalculate(k int) {
	p.dvelocity[k].X = p.dMomentum[k].X * p.massInverse
	p.dvelocity[k].Y = p.dMomentum[k].Y * p.massInverse
}

type Space struct {
	time float64
	Particles []Particle
}

func MkSpace() *Space {
	return new(Space)
}

func (s *Space) Particle(id int) *Particle {
	return &s.Particles[id]
}

func (s *Space) MkParticle(mass float64) *Particle {
	id := len(s.Particles)
	p := Particle{Id: id}
	p.SetMass(mass)
	s.Particles = append(s.Particles, p)
	return &s.Particles[id]
}

func (s *Space) Step(dt float64) {
	evaluate(s.Particles, 0, 0)
	evaluate(s.Particles, dt/2, 1)
	evaluate(s.Particles, dt/2, 2)
	evaluate(s.Particles, dt, 3)
	updateParameters(s.Particles)
	s.time += dt
}

const G = 6.67384
var _ = math.MaxFloat64 // TODO: remove this

func evaluate(particles []Particle, dt float64, k int) {
	if k == 0 {
		for i := range particles {
			p := &particles[i]

			for i := 0; i < 4; i++ {
				p.dPosition[i] = p.Position
				p.dMomentum[i] = p.Momentum
				p.dRecalculate(i)
				p.dForce[i] = vect.V{}
			}
		}
	}
	for i := range particles {
		p := &particles[i]
		for j := i + 1; j < len(particles); j++ {
			q := &particles[j]

			dx := p.dPosition[k].X - q.dPosition[k].X
			dy := p.dPosition[k].Y - q.dPosition[k].Y
			if dx < 0.01 && dy < 0.01 {
				continue
			}

			dSquared := dx*dx + dy*dy
			distance := math.Sqrt(dSquared)
			mag := dSquared * distance

			fX := G*p.Mass*q.Mass*dx / mag
			fY := G*p.Mass*q.Mass*dy / mag

			p.dForce[k].X -= fX
			p.dForce[k].Y -= fY

			q.dForce[k].X += fX
			q.dForce[k].Y += fY
		}
	}
	if k != 0 {
		for i := range particles {
			p := &particles[i]

			p.dPosition[k].X = p.Position.X + p.dvelocity[k-1].X * dt
			p.dPosition[k].Y = p.Position.Y + p.dvelocity[k-1].Y * dt

			p.dMomentum[k].X = p.Momentum.X + p.dForce[k-1].X * dt
			p.dMomentum[k].Y = p.Momentum.Y + p.dForce[k-1].Y * dt
			p.dRecalculate(k)
		}
	}
}

func updateParameters(particles []Particle) {
	for i := range particles {
		p := &particles[i]
		p.Position = calculateRK4Mean(p.dPosition)
		p.Momentum = calculateRK4Mean(p.dMomentum)
		p.Recalculate()
	}
}

func calculateRK4Mean(vs [4]vect.V) vect.V {
	const FRAC = 1.0/6.0
	return vect.V {
		X: FRAC * (vs[0].X + 2*(vs[1].X + vs[2].X) + vs[3].X),
		Y: FRAC * (vs[0].Y + 2*(vs[1].Y + vs[2].Y) + vs[3].Y),
	}
}

