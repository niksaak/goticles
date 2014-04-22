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

func mkParticle(id int) Particle {
	p := Particle{ id: id }
	p.SetMass(1)
	return p
}

func mkSpace() Space {
	s := Space{
		particles: []Particle{
		},
	}
	s.AddParticle(1)
	return s
}

func BenchmarkSpace(b *testing.B) {
	const STEP = 1.0/60
	space := mkSpace()
	particle := &space.particles[0]
	particle.ApplyImpulse(vect.V{1, 0})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Integrate(STEP)
	}
	b.Logf("t: %.4f; p: %v", space.time, particle)
}

func BenchmarkSpace1024(b *testing.B) {
	const PARTICLE_COUNT = 1024
	const STEP = 1.0/60
	space := mkSpace()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.AddParticle(1)
		particle.ApplyImpulse(vect.V{rf64(), rf64()})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Integrate(STEP)
	}
	b.Logf("space.time = %.4f", space.time)
}

