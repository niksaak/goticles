package main

import (
	"fmt"
	"github.com/go-gl/glfw3"
	"github.com/go-gl/glow/gl-core/4.4/gl"
	"github.com/go-gl/mathgl/mgl32"
	//	"github.com/niksaak/goticles"
	"github.com/niksaak/goticles/demo/engine"
	"github.com/niksaak/goticles/gputicles"
	"github.com/niksaak/goticles/vect"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

const (
	PARTICLE_MASS_DEFAULT  = 64
	PARTICLE_COUNT_DEFAULT = 16384
)

func randFloat() float64 {
	return rand.Float64() - 0.5
}

func randVect() vect.V {
	return vect.V{randFloat(), randFloat()}.Ulen().Mul(rand.Float64())
}

type MainState struct {
	space *gputicles.Space

	vao     uint32
	vbo     uint32
	program uint32

	projection        mgl32.Mat4
	projectionUniform int32
}

func (s *MainState) Init(e *engine.E) error {
	// Init Glow
	if err := gl.Init(); err != nil {
		return err
	}

	// Enable debug output
	if glfw3.ExtensionSupported("GL_ARB_debug_output") {
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS_ARB)
		gl.DebugMessageCallbackARB(gl.DebugProc(glDebugCallback), nil)
	}

	// Init shaders
	shaderProg := engine.NewShaderProgram()
	if err := shaderProg.ReadShaderFile("vert.glsl", engine.VertexShader); err != nil {
		return err
	}
	if err := shaderProg.ReadShaderFile("geom.glsl", engine.GeometryShader); err != nil {
		return err
	}
	if err := shaderProg.ReadShaderFile("frag.glsl", engine.FragmentShader); err != nil {
		return err
	}
	program, err := shaderProg.Link()
	if err != nil {
		return err
	}
	s.program = program
	gl.UseProgram(program)

	// Manage GL features
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.Enable(gl.PROGRAM_POINT_SIZE)

	// Init VAO
	gl.GenVertexArrays(1, &s.vao)
	gl.BindVertexArray(s.vao)
	e.Defer(func() {
		gl.DeleteVertexArrays(1, &s.vao)
	})

	s.projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))

	// Init VBO
	gl.GenBuffers(1, &s.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	e.Defer(func() {
		gl.DeleteBuffers(1, &s.vbo)
	})

	positionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Init space
	space, err := gputicles.New()
	if err != nil {
		return err
	}
	s.space = space
	/*
		for i := 0; i < PARTICLE_COUNT_DEFAULT; i++ {
			particle := space.MkParticle(PARTICLE_MASS_DEFAULT)
			particle.Position = randVect()
			//particle.Velocity = randVect().Div(8)
		}
	*/
	for i := 0; i < PARTICLE_COUNT_DEFAULT/2; i++ {
		particle := space.MkParticle(PARTICLE_MASS_DEFAULT)
		particle.Position = randVect().Div(2).Add(vect.V{-0.5, -0.5})
	}
	for i := PARTICLE_COUNT_DEFAULT / 2; i < PARTICLE_COUNT_DEFAULT; i++ {
		particle := space.MkParticle(PARTICLE_MASS_DEFAULT)
		particle.Position = randVect().Div(2).Add(vect.V{0.5, 0.5})
	}

	return nil
}

func (s *MainState) Update(dt float64) error {
	s.space.Step(dt)
	/*
		id := 1
		fmt.Printf(
			"position: %.5f; velocity: %.5f; acceleration %.5f\n",
			s.space.Particles[id].Position, s.space.Particles[id].Velocity,
			s.space.Particles[id].Acceleration)
	*/
	return nil
}

func (s *MainState) Render() error {
	gl.UseProgram(s.program)
	gl.UniformMatrix4fv(s.projectionUniform, 1, false, &s.projection[0])

	/*
		gl.BindVertexArray(s.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	*/
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(PARTICLE_COUNT_DEFAULT*unsafe.Sizeof([2]float32{})),
		nil,
		gl.STREAM_DRAW)
	data, err := sliceMap2Float32Buffer(gl.WRITE_ONLY, PARTICLE_COUNT_DEFAULT)
	if err != nil {
		return err
	}
	for i := range data {
		p := &s.space.Particles[i]
		if p == nil {
			return fmt.Errorf("no particle of id %d", i)
		}
		data[i][0] = float32(p.Position.X)
		data[i][1] = float32(p.Position.Y)
	}
	gl.UnmapBuffer(gl.ARRAY_BUFFER)
	data = nil

	positionAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.POINTS, 0, PARTICLE_COUNT_DEFAULT)

	return nil
}

func sliceMap2Float32Buffer(mode uint32, length int) ([][2]float32, error) {
	data := gl.MapBuffer(gl.ARRAY_BUFFER, mode)
	if data == nil {
		return nil, fmt.Errorf("buffer data pointer is nil")
	}
	header := reflect.SliceHeader{
		Data: uintptr(data),
		Len:  length,
		Cap:  length,
	}
	return *((*[][2]float32)(unsafe.Pointer(&header))), nil
}

func (s *MainState) Resize(width, height int) {
	aspect := float32(width) / float32(height)
	if width > height {
		s.projection = mgl32.Ortho2D(-aspect, aspect, -1, 1)
	} else {
		s.projection = mgl32.Ortho2D(-1, 1, -1/aspect, 1/aspect)
	}
	gl.Viewport(0, 0, int32(width), int32(height))
}

func glDebugCallback(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {
	fmt.Fprintf(
		os.Stderr,
		"Debug (source: %d, type: %d severity: %d): %s\n",
		source, gltype, severity, message)
}

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	params := engine.Params{
		Version: [2]int{4, 4},
		Size:    [2]int{512, 512},
	}
	engine := engine.New()
	if err := engine.Initialize("GPU TICKLES", params, &MainState{}); err != nil {
		panic(err)
	}
	if err := engine.Run(); err != nil {
		panic(err)
	}
	if err := engine.Shutdown(); err != nil {
		panic(err)
	}
}
