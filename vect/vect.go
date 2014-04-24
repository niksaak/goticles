// Package Vect provides basic vector operations.
package vect

import (
	"math"
)

// The V type designates a two-dimensional vector.
type V struct{ X, Y float64 }

// Eql returns true when vectors v and a are equal.
func (v V) Eql(a V) bool {
	return v == a
}

// Add a to v.
func (v V) Add(a V) V {
	return V{v.X + a.X, v.Y + a.Y}
}

// Subtract a from v.
func (v V) Sub(a V) V {
	return V{v.X - a.X, v.Y - a.Y}
}

// Multiply v by m.
func (v V) Mul(m float64) V {
	return V{v.X * m, v.Y * m}
}

// Divide v bY d.
func (v V) Div(d float64) V {
	return V{v.X / d, v.Y / d}
}

// Dot product of vectors v and m.
func (v V) Dot(m V) float64 {
	return (v.X * m.X) + (v.Y * m.Y)
}

// Dst returns the distance between vectors v and d.
func (v V) Dst(d V) float64 {
	return v.Sub(d).Len()
}

// Dst returns the squared distance between vectors v and d.
func (v V) DstSq(d V) float64 {
	return v.Sub(d).LenSq()
}

// Negate a vector.
func (v V) Neg() V {
	return V{-v.X, -v.Y}
}

// Len returns the length of a vector.
func (v V) Len() float64 {
	return math.Sqrt(v.Dot(v))
}

// LenSq returns the length of a vector, squared.
func (v V) LenSq() float64 {
	return v.Dot(v)
}

