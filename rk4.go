package goticles

import (
	"github.com/niksaak/goticles/vect"
)

type State struct {
	position    vect.V
	momentum    vect.V
	force       vect.V
	velocity    vect.V
	mass        float64
	massInverse float64
}

func (s State) Position() vect.V {
	return s.position
}

func (s *State) SetPosition(position vect.V) {
	s.position = position
}

func (s State) Momentum() vect.V {
	return s.momentum
}

func (s *State) SetMomentum(momentum vect.V) {
	s.momentum = momentum
	s.recalculate()
}

func (s State) Force() vect.V {
	return s.force
}

func (s *State) SetForce(force vect.V) {
	s.force = force
}

func (s *State) ApplyImpulse(impulse vect.V) {
	s.momentum = s.momentum.Add(impulse)
	s.recalculate()
}

func (s State) Mass() float64 {
	return s.mass
}

func (s State) MassInverse() float64 {
	return s.massInverse
}

func (s *State) SetMass(mass float64) {
	s.mass = mass
	s.massInverse = 1 / mass
	s.recalculate()
}

func (s State) Velocity() vect.V {
	return s.velocity
}

func (s *State) recalculate() {
	s.velocity = s.momentum.Mul(s.massInverse)
}

type Derivative struct {
	momentum vect.V
	force    vect.V
}

func evaluate(initial State, t, dt float64, d Derivative) Derivative {
	var s State
	s.position = initial.position.Add(d.momentum.Mul(dt))
	s.momentum = initial.momentum.Add(d.force.Mul(dt))

	var ret Derivative
	ret.momentum = s.momentum
	ret.force = initial.force
	return ret
}

func integrate(s *State, t, dt float64) {
	a := evaluate(*s, t, 0, Derivative{})
	b := evaluate(*s, t, dt/2, a)
	c := evaluate(*s, t, dt/2, b)
	d := evaluate(*s, t, dt, c)

	dMomentum := a.momentum.Add(b.momentum.Add(c.momentum).Mul(2).Add(d.momentum)).Mul(1.0 / 6)
	dForce := a.force.Add(b.force.Add(c.force).Mul(2).Add(d.force)).Mul(1.0 / 6)

	s.position = s.position.Add(dMomentum.Mul(dt))
	s.momentum = s.momentum.Add(dForce.Mul(dt))
	s.recalculate()
	s.force = vect.V{0, 0}
}

