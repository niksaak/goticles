package goticles

import (
	"github.com/niksaak/goticles/vect"
	"math/rand"
	"testing"
)

// rf64 returns random float64 in range [-0.5, 0.5).
func rf64() float64 {
	return rand.Float64() - 0.5
}

func TestOneBodyPositionChange(t *testing.T) {
	const STEP = 1.0/1
	SuperVerbose = true
	space := MkSpace()
	if count := len(space.Particles); count != 1 {
		t.Errorf("particle count is not 1 but %d", count)
	}
	p := space.MkParticle(1)
	p.SetPosition(vect.V{0, 0})
	p.SetVelocity(vect.V{1, 0})
	t.Logf("%.4f: %v", 0.0, p)
	/*
	for c := 0.0; c < 1; c += STEP {
		space.Step(STEP)
		t.Logf("%.4f: %v", c, p)
	}
	*/
	c := 0.0
	for {
		space.Step(STEP)
		t.Logf("%.4f: %v", c, p)
		if c > 1 {
			break
		}
		c += STEP
	}
	SuperVerbose = false
}

func TestTwoBodyPositionChanges(t *testing.T) {
	const STEP = 1.0/60
	space := MkSpace()
	p1 := space.MkParticle(1)
	p1.SetPosition(vect.V{0, 0})
	p2 := space.MkParticle(1)
	p2.SetPosition(vect.V{1, 0})
	t.Log("Before:")
	t.Log(p1)
	t.Log(p2)
	space.Step(STEP)
	t.Log("After:")
	t.Log(p1)
	t.Log(p2)
}

func BenchmarkSpace(b *testing.B) {
	const STEP = 1.0/60
	space := MkSpace()
	particle := space.MkParticle(1)
	particle.SetVelocity(vect.V{1, 0})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.time, particle)
}

func BenchmarkSpace512(b *testing.B) {
	const PARTICLE_COUNT = 512
	const STEP = 1.0/60
	space := MkSpace()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.MkParticle(1)
		particle.SetVelocity(vect.V{rf64(), rf64()})
	}
	particle := &space.Particles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.time, particle)
}
func BenchmarkSpace1024(b *testing.B) {
	const PARTICLE_COUNT = 1024
	const STEP = 1.0/60
	space := MkSpace()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.MkParticle(1)
		particle.SetVelocity(vect.V{rf64(), rf64()})
	}
	particle := &space.Particles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.time, particle)
}

