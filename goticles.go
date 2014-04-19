package goticles

import (
	//"github.com/go-gl/gl"
	//"github.com/go-gl/glfw3"
	"fmt"
)

type Vect struct{ x, y float64 }

// X returns the X coordinate of a vector.
func (v Vect) X() float64 {
	return v.x
}

// Y returns the Y coordinate of the vector.
func (v Vect) Y() float64 {
	return v.y
}

// Add a to v.
func (v Vect) Add(a Vect) Vect {
	return Vect{v.x + a.x, v.y + a.y}
}

// Subtract a from v.
func (v Vect) Sub(a Vect) Vect {
	return Vect{v.x - a.x, v.y - a.y}
}

// Multiply v by m.
func (v Vect) Mul(m float64) Vect {
	return Vect{v.x * m, v.y * m}
}

// Divide v by d.
func (v Vect) Div(d float64) Vect {
	return Vect{v.x / d, v.y / d}
}

// Negate a vector.
func (v Vect) Neg() Vect {
	return Vect{-v.x, -v.y}
}

type state struct {
	position Vect
	velocity Vect
}

/*
type derivative state

type Particle struct {
	s state
	d derivative
}

func evaluate(
	initial state,
	t, dt float64,
	d derivative,
	accel func(s state, t float64),
) derivative {
	if accel == nil {
		panic("No accel func")
	}
	var s state
	s.position = initial.position.Add(d.position.Mul(dt))
	s.velocity = initial.velocity.Add(d.velocity.Mul(dt))

	var ret derivative
	ret.position = state.velocity
	ret.velocity = accel(s, t+dt)
	return ret
}
*/

type State struct {
	position    Vect
	momentum    Vect
	velocity    Vect
	mass        float64
	massInverse float64
}

func (s State) Position() Vect {
	return s.position
}

func (s *State) SetPosition(position Vect) {
	s.position = position
}

func (s State) Momentum() Vect {
	return s.momentum
}

func (s *State) SetMomentum(momentum Vect) {
	s.momentum = momentum
	s.recalculate()
}

func (s *State) ApplyImpulse(impulse Vect) {
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
	s.massInverse = 1/mass
	s.recalculate()
}

func (s State) Velocity() Vect {
	return s.velocity
}

func (s *State) recalculate() {
	s.velocity = s.momentum.Mul(s.massInverse)
}

type Derivative struct {
	velocity Vect
	force    Vect
}

type ForceFunc func(s State, t float64) Vect

func evaluate(initial State, t, dt float64, d Derivative, accel ForceFunc) Derivative {
	var s State
	s.position = initial.position.Add(d.velocity.Mul(dt))
	s.velocity = initial.velocity.Add(d.force.Mul(dt))

	var ret Derivative
	ret.velocity = s.velocity
	ret.force = accel(s, t+dt)
	return ret
}

func integrate(s *State, t, dt float64, force ForceFunc) {
	a := evaluate(*s, t, 0, Derivative{}, force)
	b := evaluate(*s, t, dt/2, a, force)
	c := evaluate(*s, t, dt/2, b, force)
	d := evaluate(*s, t, dt, c, force)

	dVelocity := a.velocity.Add(b.velocity.Add(c.velocity).Mul(2).Add(d.velocity)).Mul(1.0/6)
	dForce := a.force.Add(b.force.Add(c.force).Mul(2).Add(d.force)).Mul(1.0/6)

	s.position = s.position.Add(dVelocity.Mul(dt))
	s.velocity = s.velocity.Add(dForce.Mul(dt))
	s.recalculate()
}

type Particle struct {
	id int64
	State
}

func (p *Particle) String() string {
	return fmt.Sprintf(
		"#%d @(%.4f, %.4f) V(%.4f, %.4f) ->(%.4f, %.4f)",
		p.id,
		p.State.position.x,
		p.State.position.y,
		p.State.velocity.x,
		p.State.velocity.y,
		p.State.momentum.x,
		p.State.momentum.y,
	)
}

type Space struct {
	time float64
	particles []Particle
	force ForceFunc
}

func (s *Space) Integrate(dt float64) {
	force := s.force
	for i := range s.particles {
		integrate(&s.particles[i].State, s.time, dt, force)
	}
	s.time += dt
}

