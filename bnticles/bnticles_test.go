package bnticles

import(
	"math/rand"
	"fmt"
	"github.com/niksaak/goticles/vect"
	"testing"
)

const TREESTRING_INDENT = "  "

func treeString(nd node, prefix string) (s string) {
	s = prefix
	switch n := nd.(type) {
	case *quad:
		s += fmt.Sprintf("quad %v; COM: %v, m: %v\n", n.BB, n.position, n.mass)
		for _, child := range n.children {
			if child == nil {
				s += prefix + TREESTRING_INDENT + "nil\n"
			} else {
				s += treeString(child, prefix + TREESTRING_INDENT)
			}
		}
	case *particle:
		s += fmt.Sprintf("particle @%v; ->%v; m: %v\n",
			n.position, n.velocity, n.mass)
	}
	return s
}

// rf64 returns random float64 in range [-0.5, 0.5).
func rf64() float64 {
	return rand.Float64() - 0.5
}

func spaceTree(s *Space) *quad {
	s.evaluate1()
	qtree := NewQuad(BB{{-1, 1}, {1, -1}})
	qtree.insertSlice(s.bnParticles[0])
	return qtree
}

func makeSpace(count int, tb testing.TB) *Space {
	space := New()
	for i := 0; i < count; i++ {
		p := space.MkParticle(1)
		p.Position = vect.V{rf64(), rf64()}
	}
	if result := len(space.Particles); count != result {
		tb.Fatalf("particle count is not %d but %d", count, result)
	}
	return space
}

func TestQuadTreeBuilding(t *testing.T) {
	const COUNT = 8
	space := makeSpace(COUNT, t)
	quad := spaceTree(space)
	for _, p := range space.bnParticles[0] {
		t.Logf("created %s", treeString(&p, ""))
	}
	t.Logf("\n%s", treeString(quad, ""))
}

func (p *Particle) String() string {
	return fmt.Sprintf("particle @%v, ->%v, m: %v",
		p.Position, p.Velocity, p.Mass)
}

func logParticles(s *Space, tb testing.TB) {
	for _, p := range s.Particles {
		tb.Log(p.String())
	}
}

func TestOneBodyPositionChanges(t *testing.T) {
	const STEP = 1.0 / 10
	space := makeSpace(1, t)
	space.Particles[0].Position = vect.V{0, 0}
	space.Particles[0].Velocity = vect.V{1, 1}
	logParticles(space, t)
	for time := float64(0); time < 1; time += STEP {
		space.Step(STEP)
		logParticles(space, t)
	}
}

func TestTwoBodyPositionChanges(t *testing.T) {
	const STEP = 1.0 / 10
	space := makeSpace(2, t)
	logParticles(space, t)
	for time := float64(0); time < 1; time += STEP {
		space.Step(STEP)
		logParticles(space, t)
	}
}

func BenchmarkQuadTreeBuilding(b *testing.B) {
	const COUNT = 1024
	space := New()
	// do some dirty deeds to purify the bench results
	for i := 0; i < COUNT; i++ {
		p := space.MkParticle(1)
		p.Position = vect.V{rf64(), rf64()}
	}
	space.evaluate1()
	particles := space.bnParticles[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quad := NewQuad(BB{{-1, 1}, {1, -1}})
		quad.insertSlice(particles)
	}
}

func BenchmarkBarnesHutSimulation(b *testing.B) {
	const STEP = 1.0 / 10
	const COUNT = 1024
	space := makeSpace(COUNT, b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
}
