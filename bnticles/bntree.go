package bnticles

import (
	"fmt"
	"github.com/niksaak/goticles/vect"
)

type node interface {
	Position() vect.V
	Mass() float64
	Children() ([4]node, bool)
}

func nodeLen(n node) int {
	if n == nil {
		return 0
	}
	if children, ok := n.Children(); !ok {
		return 1
	} else {
		var count int
		for _, child := range children {
			count += nodeLen(child)
		}
		return count
	}
}

type quad struct {
	BB
	position vect.V // center of mass
	mass     float64
	children [4]node
}

var _ node = &quad{}

func NewQuad(bb BB) *quad {
	return &quad{BB: bb}
}

func (q *quad) Position() vect.V { return q.position }

func (q *quad) Mass() float64 { return q.mass }

func (q *quad) Children() ([4]node, bool) { return q.children, false }

func (q *quad) recalculate() {
	var (
		pos  vect.V
		mass float64
	)
	for _, n := range q.children {
		if n == nil {
			continue
		}
		pos = pos.Add(n.Position().Mul(n.Mass()))
		mass += n.Mass()
	}
	q.position = pos.Div(mass)
	q.mass = mass
}

func (q *quad) sideFor(v vect.V) (side Side) {
	center := q.Center()
	if v.X < center.X {
		side |= west
	}
	if v.Y < center.Y {
		side |= south
	}
	return side
}

func (q *quad) insert(p *particle) {
	if q == nil {
		panic("insert: quad is nil")
	}
	side := q.sideFor(p.position)
	if q.children[side] == nil {
		q.children[side] = p
		q.recalculate()
		return
	}

	switch n := q.children[side].(type) {
	case *quad:
		n.insert(p)
	case *particle:
		q.children[side] = NewQuad(q.Side(q.sideFor(p.position)))
		q.insert(p)
		q.insert(n)
	}
	q.recalculate()
}

func (q *quad) insertSlice(particles []particle) {
	for i := range particles {
		q.insert(&particles[i])
	}
}

type particle struct {
	position vect.V
	velocity vect.V
	mass     float64
}

var _ node = &particle{}

func (p *particle) Position() vect.V { return p.position }

func (p *particle) Mass() float64 { return p.mass }

func (p *particle) Children() ([4]node, bool) { return [4]node{}, false }

const (
	THETA = 0.6
	G     = 6.67384e-11
)

func (p *particle) force(n node) vect.V {
	if n == nil {
		return vect.V{}
	}
	distV := p.position.Sub(n.Position())
	distSq := distV.LenSq()
	distU := distV.Ulen()
	return distU.Mul(G * p.mass * n.Mass() / distSq)
}

func (p *particle) treeForce(tree node) vect.V {
	if tree == nil {
		return vect.V{}
	}
	switch n := tree.(type) {
	case *quad:
		dist := p.position.Dst(n.position)
		size := n.Size()
		if size/dist < THETA {
			return p.force(n)
		} else {
			force := vect.V{}
			for _, child := range n.children {
				force = force.Add(p.treeForce(child))
			}
			return force
		}
	case *particle:
		if n == p {
			return vect.V{}
		}
		return p.force(n)
	default:
		panic(fmt.Errorf("treeForce: bad argument type - %T", tree))
	}
}
