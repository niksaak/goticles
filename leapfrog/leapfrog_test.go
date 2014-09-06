package leapfrog

import (
	"github.com/niksaak/goticles/vect"
	"math/rand"
	"testing"
)

func randVect() vect.V {
	return vect.V{rand.Float64() - 0.5, rand.Float64() - 0.5}
}

func uRandVect() vect.V {
	return randVect().Ulen()
}

func makeSpace(particleCount int) *Space {
	s := New()
	for i := 0; i < particleCount; i++ {
		p := s.MkParticle(1)
		p.Position = uRandVect()
	}
	return s
}

/*
func TestOneBodyIntegration(t *testing.T) {
	const dt = 1.0/10.0
	space := New()
	particle := space.MkParticle(1)
	particle.Velocity = vect.V{1, 0}
	position := particle.Position
	t.Logf("t: 0.00, position: %v", particle.Position)
	for time := float64(0); time < 1; time += dt {
		space.Step(dt)
		t.Logf("t: %.2f; position: %v", time, particle.Position)
	}
	newPosition := particle.Position
	if position.Eql(newPosition) {
		t.Errorf(
			"positions %v and %v are equal after %.0f iterations",
			position, newPosition, 1.0 / dt)
	}
}
*/

func TestTwoBodiesIntegration(t *testing.T) {
	const dt = 1.0/10.0
	space := makeSpace(2)
	position := space.Particle(0).Position
	for time := float64(0); time < 1; time += dt {
		space.Step(dt)
	}
	newPosition := space.Particle(0).Position
	if position.Eql(newPosition) {
		t.Errorf(
			"positions %v and %v are equal after %.0f iterations",
			position, newPosition, 1.0 / dt)
	}
}

func BenchmarkSimulation2(b *testing.B) {
	const dt = 1.0/100
	space := makeSpace(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(dt)
	}
}

func BenchmarkSimulation32(b *testing.B) {
	const dt = 1.0/100
	space := makeSpace(32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(dt)
	}
}

func BenchmarkSimulation256(b *testing.B) {
	const dt = 1.0/100
	space := makeSpace(256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(dt)
	}
}

func BenchmarkSimulation1024(b *testing.B) {
	const dt = 1.0/100
	space := makeSpace(1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		space.Step(dt)
	}
}
