package rk4

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
	const STEP = 1.0 / 100
	space := MkSpaceRK4()
	p := space.MkParticle(1)
	if count := len(space.Particles); count != 1 {
		t.Errorf("particle count is not 1 but %d", count)
	}
	p.SetPosition(vect.V{0, 0})
	p.SetVelocity(vect.V{1, 1})
	t.Logf("%.4f: %v", space.Time, p)
	for c := 0.0; c < 1; c += STEP {
		space.Step(STEP)
		t.Logf("%.4f: %v", space.Time, p)
	}
}

func TestTwoBodyPositionChanges(t *testing.T) {
	const STEP = 1.0 / 60
	const MASS = 10000
	space := MkSpaceRK4()
	p1 := space.MkParticle(MASS)
	p1.SetPosition(vect.V{0, 0})
	p2 := space.MkParticle(MASS)
	p2.SetPosition(vect.V{0.5, 0.5})
	t.Log("Before:")
	t.Log(p1)
	t.Log(p2)
	space.Step(STEP)
	t.Log("After:")
	t.Log(p1)
	t.Log(p2)
}

func BenchmarkSpaceRK41(b *testing.B) {
	const STEP = 1.0 / 60
	const MASS = 10000
	space := MkSpaceRK4()
	particle := space.MkParticle(MASS)
	particle.SetVelocity(vect.V{1, 0})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.Time, particle)
}

func BenchmarkSpaceRK4512(b *testing.B) {
	const PARTICLE_COUNT = 512
	const MASS = 10000
	const STEP = 1.0 / 60
	space := MkSpaceRK4()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.MkParticle(MASS)
		particle.SetVelocity(vect.V{rf64(), rf64()})
	}
	particle := &space.Particles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.Time, particle)
}

func BenchmarkSpaceRK41024(b *testing.B) {
	const PARTICLE_COUNT = 1024
	const STEP = 1.0 / 60
	const MASS = 10000
	space := MkSpaceRK4()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.MkParticle(MASS)
		particle.SetVelocity(vect.V{rf64(), rf64()})
	}
	particle := &space.Particles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.Time, particle)
}
