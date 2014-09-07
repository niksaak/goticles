package gputicles

import (
	"fmt"
	"github.com/go-gl/glow/gl-core/4.4/gl"
	"github.com/niksaak/goticles"
	"github.com/niksaak/goticles/demo/engine"
	"github.com/niksaak/goticles/vect"
	"reflect"
	"unsafe"
)

const (
	WORK_GROUP_LOCAL_SIZE = 256
	TILED_SHADER = true
)

type Space struct {
	Particles []goticles.P

	accelProgram uint32
	vao          uint32
	vbo          uint32
}

type gpuParticle struct {
	position     [2]float32
	acceleration [2]float32
	mass         float32
	_            [3]float32
}

func New() (s *Space, err error) {
	s = new(Space)
	s.accelProgram, err = loadShader()
	if err != nil {
		return nil, err
	}
	gl.GenVertexArrays(1, &s.vao)
	gl.GenBuffers(1, &s.vbo)

	gl.UseProgram(s.accelProgram)
	gl.BindVertexArray(s.vao)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, s.vbo)

	return s, nil
}

func (s *Space) MkParticle(mass float64) *goticles.P {
	id := len(s.Particles)
	s.Particles = append(s.Particles, goticles.P{
		Id:   id,
		Mass: mass,
	})
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
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, s.vbo)

	gl.BufferData(
		gl.SHADER_STORAGE_BUFFER,
		len(s.Particles)*int(unsafe.Sizeof(gpuParticle{})),
		nil,
		gl.DYNAMIC_COPY)
	{ // send
		mappedParticles, err := sliceMapParticlesBuffer(gl.WRITE_ONLY, len(s.Particles))
		if err != nil {
			panic(err)
		}
		for i, p := range s.Particles {
			mappedParticles[i] = gpuParticle{
				position: [2]float32{float32(p.Position.X), float32(p.Position.Y)},
				mass:     float32(p.Mass),
			}
		}
		gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	}
	gl.DispatchCompute(uint32(len(s.Particles))/WORK_GROUP_LOCAL_SIZE, 1, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	{ // receive
		mappedParticles, err := sliceMapParticlesBuffer(gl.READ_ONLY, len(s.Particles))
		if err != nil {
			panic(err)
		}
		for i, q := range mappedParticles {
			p := &s.Particles[i]
			p.Position = p.Position.Add(p.Velocity.Mul(dt)).Add(
				p.Acceleration.Mul(dt * dt).Mul(0.5))
			oldAcceleration := p.Acceleration
			p.Acceleration = vect.V{
				float64(q.acceleration[0]), float64(q.acceleration[1]),
			}
			p.Velocity = p.Velocity.Add(
				p.Acceleration.Add(oldAcceleration).Mul(0.5).Mul(dt))
		}
		gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	}

	gl.UseProgram(uint32(prevProgram))
}

func sliceMapParticlesBuffer(mode uint32, length int) ([]gpuParticle, error) {
	data := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, mode)
	if data == nil {
		return nil, fmt.Errorf("buffer data pointer is nil")
	}
	header := reflect.SliceHeader{
		Data: uintptr(data),
		Len:  length,
		Cap:  length,
	}
	slice := (*[]gpuParticle)(unsafe.Pointer(&header))
	return *slice, nil
}

func loadShader() (shader uint32, err error) {
	accelProgram := engine.NewShaderProgram()
	if !TILED_SHADER {
		err = accelProgram.ReadShaderFile("accelerator.comp.glsl", engine.ComputeShader)
	} else {
		err = accelProgram.ReadShaderFile("accelerator_tiled.comp.glsl", engine.ComputeShader)
	}
	if err != nil {
		return 0, err
	}
	return accelProgram.Link()
}
