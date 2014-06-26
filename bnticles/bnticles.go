package bnticles

import (
	"fmt"
	"github.com/niksaak/goticles/vect"
)

// The Spacer interface is satisfied by particle simulation types.
type Spacer interface {
	Particle(id int) *Particle         // get particle
	MkParticle(mass float64) *Particle // add particle to the simulation
	RmParticle(p *Particle)            // remove particle from the simulation
	Step(dt float64)                   // step simulation by dt
}

// The Particle represents a unit in a simulation.
type Particle struct {
	Id          int // identifier
	Position    vect.V
	Velocity    vect.V
	Force       vect.V
	Mass        float64
	massInverse float64 // = (1 / Particle.Mass)
}

// The values of those constants correspond to positions in children array.
const (
	ne = iota
	se
	sw
	nw
)

// The bnBB type represents an axis-aligned bounding box, storing its
// NorthWest and SouthEast corners.
type bnBB [2]vect.V

func (b bnBB) Center() vect.V {
	return vect.V{
		X: 0.5 * (b[0].X + b[1].X),
		Y: 0.5 * (b[0].Y + b[1].Y),
	}
}

func (b bnBB) Size() float64 {
	return b[1].X - b[0].X
}

func (b bnBB) Query(v vect.V) bool {
	return v.X >= b[0].X && v.Y >= b[0].Y &&
		v.X <= b[1].X && v.Y <= b[1].Y
}

func (b bnBB) SubBB(side int) bnBB {
	center := b.Center()
	switch side {
	case ne:
		return bnBB{
			vect.V{center.X, b[0].Y},
			vect.V{b[1].X, center.Y},
		}
	case se:
		return bnBB{
			center,
			b[1],
		}
	case sw:
		return bnBB{
			vect.V{b[0].X, center.Y},
			vect.V{center.X, b[1].Y},
		}
	case nw:
		return bnBB{
			b[0],
			center,
		}
	default:
		return bnBB{}
	}
}

// A bnNode is a node of Barnes-Hut tree, which can be a particle or a subtree.
type bnNode interface {
	Position() vect.V
	Mass() float64
	Children() *[4]bnNode
}

func bnLen(n bnNode) int {
	if n == nil {
		return 0
	}
	if children := n.Children(); children == nil {
		return 1
	} else {
		var count int
		for _, child := range children {
			count += bnLen(child)
		}
		return count
	}
}

type bnQuad struct {
	bnBB
	com      vect.V
	mass     float64
	children [4]bnNode
}

func bnMkQuad(bb bnBB) *bnQuad {
	q := &bnQuad{bnBB: bb}
	return q
}

func (t *bnQuad) Position() vect.V { return t.com }

func (t *bnQuad) Mass() float64 { return t.mass }

func (t *bnQuad) Children() *[4]bnNode { return &t.children }

func (t *bnQuad) updateParameters() {
	var (
		pos vect.V
		mass float64
	)
	for _, node := range t.children {
		if node == nil {
			continue
		}
		pos = pos.Add(node.Position().Mul(node.Mass()))
		mass += node.Mass()
	}
	t.com = pos.Div(mass)
	t.mass = mass
}

func (t *bnQuad) side(v vect.V) int {
	center := vect.V{
		X: 0.5 * (t.bnBB[0].X + t.bnBB[1].X),
		Y: 0.5 * (t.bnBB[0].Y + t.bnBB[1].Y),
	}
	if v.X < center.X { // west
		if v.Y < center.Y {
			return sw // south
		} else {
			return nw // north
		}
	} else { // east
		if v.Y < center.Y {
			return se // south
		} else {
			return ne // north
		}
	}
}

func (t *bnQuad) insert(p *bnParticle) {
	side := t.side(p.position)
	if t.children[side] == nil {
		t.children[side] = p
		t.updateParameters()
		return
	}

	switch n := t.children[side].(type) {
	case *bnQuad:
		n.insert(p)
	case *bnParticle:
		t.children[side] = bnMkQuad(t.bnBB.SubBB(side))
		t.children[side].insert(p)
		t.children[side].insert(n)
	}
	t.updateParameters()
	return
}

type bnParticle struct {
	position vect.V
	velocity vect.V
	mass     float64
}

func (p bnParticle) Position() vect.V { return p.position }

func (p bnParticle) Mass() float64 { return p.mass }

func (p bnParticle) Children() *[4]bnNode { return nil }

const (
	THETA = 0.5
	G = 6.67384e-11
)

func (p *bnParticle) force(n bnNode) vect.V {
	distV := p.position.Sub(n.Position())
	distSq := distV.LenSq()
	distU := distV.Ulen()
	return distU.Mul(G * p.mass * n.Mass() / distSq)
}

func (p *bnParticle) treeForce(node bnNode) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *bnQuad:
		dist := p.position.Dst(node.Position())
		size := n.Size()
		if size / dist < THETA {
			p.velocity = p.velocity.Sub(p.force(n))
		} else {
			for _, child := range n.children {
				p.treeForce(child, target)
			}
		}
	case *bnParticle:
		p.velocity = p.velocity.Sub(p.force(n))
	}
}

type BNSpace struct {
	particles   []Particle
	bnParticles [4][]*bnParticle
}

func MkBNSpace() *BNSpace {
	s := new(BNSpace)
	s.particles = []Particle{}
	s.bnParticles = [4][]*bnParticle{}
	for i := range s.bnParticles {
		s.bnParticles[i] = []bnParticle{}
	}
	return s
}

func (s *BNSpace) Particle(id int) *Particle { return &s.particles[id] }

func (s *BNSpace) MkParticle(mass float64) *Particle {
	id := len(s.particles)
	s.particles = append(s.particles, Particle{
		Id: id,
		Mass: mass,
		massInverse: 1 / mass,
	})
	for i := range s.bnParticles {
		s.bnParticles[i] = make([]bnParticle, len(s.particles))
	}
	return &s.particles[id]
}

func (s *BNSpace) Step(dt float64) {
	s.evaluate1()
	s.evaluateK(dt, 1)
	s.evaluateK(dt, 2)
	s.evaluateK(dt, 3)
}

func (s *BNSpace) evaluate1() {
	for i, particle := range s.particles {
		s.bnParticles[0][i] = &bnParticle{
			position: particle.Position,
			velocity: particle.Velocity,
			mass:     particle.Mass,
		}
	}
}

func (s *BNSpace) evaluateK(dt float64, k int) {
	// TODO
}

