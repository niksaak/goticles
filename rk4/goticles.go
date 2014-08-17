package goticles

import (
	"fmt"
	"github.com/niksaak/goticles/vect"
)

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

const G = 6.67384e-11

