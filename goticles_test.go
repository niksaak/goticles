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

func BenchmarkSpace(b *testing.B) {
	const STEP = 1.0/60
	space := MkSpace()
	particle := space.MkParticle(1)
	particle.SetMomentum(vect.V{1, 0})
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
		particle.SetMomentum(vect.V{rf64(), rf64()})
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
		particle.SetMomentum(vect.V{rf64(), rf64()})
	}
	particle := &space.Particles[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.time, particle)
}

