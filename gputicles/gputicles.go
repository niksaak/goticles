package gputicles

import (
	"github.com/go-gl/glow/gl-core/4.4/gl"
	"github.com/niksaak/goticles/demo/engine"
	"github.com/niksaak/goticles/vect"
	"github.com/niksaak/goticles"
	"reflect"
	"unsafe"
	"fmt"
)

const (
	WORK_GROUP_LOCAL_SIZE = 256
	G                     = 6.67384e-11
)

type Space struct {
	Particles []goticles.P
	particles [4][]particle

	accelProgram uint32
	vao          uint32
	vbo          uint32
}

type particle struct {
	position     vect.V
	velocity     vect.V
	acceleration vect.V
	mass         float64
}

func New() (s *Space, err error) {
	s.accelProgram, err = loadShader()
	if err != nil {
		return nil, err
	}

	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &s.vbo)

	s = new(Space)
	return s, nil
}

func (s *Space) MkParticle(mass float64) *goticles.P {
	id := len(s.Particles)
	s.Particles = append(s.Particles, goticles.P{
		Id:   id,
		Mass: mass,
	})
	for i := range s.particles {
		s.particles[i] = append(s.particles[i], particle{})
	}
	return &s.Particles[id]
}

func (s *Space) RmParticle(p *goticles.P) {
	// TODO
}

func (s *Space) Step(dt float64) {
	var prevProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &prevProgram)
	gl.UseProgram(s.accelProgram)
	gl.BindVertexArray(s.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)

	s.evaluate1()
	s.evaluateK(dt/2, 1)
	s.evaluateK(dt/2, 2)
	s.evaluateK(dt, 3)
	s.applyState()

	gl.UseProgram(uint32(prevProgram))
}

func (s *Space) evaluate1() {
	for i, p := range s.Particles {
		s.particles[0][i] = particle{
			position: p.Position,
			velocity: p.Velocity,
			mass:     p.Mass,
		}
	}
}

func (s *Space) evaluateK(dt float64, k int) {
	count := int(uintptr(len(s.particles[0]))*unsafe.Sizeof(particle{}))
	gl.BufferData(
		gl.ARRAY_BUFFER,
		count,
		nil,
		gl.DYNAMIC_COPY)
	{
		mappedParticles, err := sliceMapParticlesBuffer(gl.WRITE_ONLY, count)
		if err != nil {
			panic(err)
		}
		copy(mappedParticles, s.particles[k-1])
		gl.UnmapBuffer(gl.ARRAY_BUFFER)
	}
	gl.DispatchCompute(uint32(count / WORK_GROUP_LOCAL_SIZE), 1, 1) // accelerations
	gl.Finish()
	{
		mappedParticles, err := sliceMapParticlesBuffer(gl.READ_ONLY, count)
		if err != nil {
			panic(err)
		}
		copy(s.particles[k], mappedParticles)
		gl.UnmapBuffer(gl.ARRAY_BUFFER)
	}
	for i := range s.particles[k] {
		p := &s.particles[k][i]
		prev := &s.particles[k-1][i]
		start := &s.particles[0][i]
		p.position = start.position.Add(prev.velocity.Mul(dt))
		p.velocity = start.velocity.Add(prev.acceleration.Mul(dt))
	}
}

func sliceMapParticlesBuffer(mode uint32, length int) ([]particle, error) {
	data := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)
	if data == nil {
		return nil, fmt.Errorf("buffer data pointer is nil")
	}
	header := reflect.SliceHeader{
		Data: uintptr(data),
		Len:  length,
		Cap:  length,
	}
	slice := (*[]particle)(unsafe.Pointer(&header))
	return *slice, nil
}

func (s *Space) applyState() {
	for i := range s.Particles {
		p := &s.Particles[i]
		p.Position = rk4mean(
			s.particles[0][i].position,
			s.particles[1][i].position,
			s.particles[2][i].position,
			s.particles[3][i].position)
		p.Velocity = rk4mean(
			s.particles[0][i].velocity,
			s.particles[1][i].velocity,
			s.particles[2][i].velocity,
			s.particles[3][i].velocity)
		p.Acceleration = rk4mean(
			s.particles[0][i].acceleration,
			s.particles[1][i].acceleration,
			s.particles[2][i].acceleration,
			s.particles[3][i].acceleration)
	}
}

func rk4mean(k1, k2, k3, k4 vect.V) vect.V {
	return vect.V{
		X: (1.0 / 6.0) * (k1.X + 2*(k2.X+k3.X) + k4.X),
		Y: (1.0 / 6.0) * (k1.Y + 2*(k2.Y+k3.Y) + k4.Y),
	}
}

func loadShader() (shader uint32, err error) {
	accelProgram := engine.NewShaderProgram()
	err = accelProgram.ReadShaderFile("accelerator.comp.glsl", engine.ComputeShader)
	if err != nil {
		return 0, err
	}
	return accelProgram.Link()
}
