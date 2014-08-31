package main

import (
	"fmt"
	"github.com/go-gl/glfw3"
	"github.com/go-gl/glow/gl-core/4.4/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/niksaak/goticles"
	"github.com/niksaak/goticles/bnticles"
	"github.com/niksaak/goticles/demo/engine"
	//"github.com/niksaak/goticles/rk4"
	"github.com/niksaak/goticles/vect"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"unsafe"
)

func randFloat() float64 {
	return rand.Float64() - 0.5
}

func randVect() vect.V {
	return vect.V{randFloat(), randFloat()}.Ulen().Mul(rand.Float64())
}

type Spacer interface {
	Particle(id int) *goticles.P
	MkParticle(mass float64) *goticles.P
	Step(dt float64)
}

const (
	PARTICLE_MASS_DEFAULT  = 96
	PARTICLE_COUNT_DEFAULT = 3072
)

type MainState struct {
	space        Spacer
	vertexArray  uint32
	vertexBuffer uint32
	program      uint32
}

var projection = mgl32.Ident4()

func (s *MainState) Init(eng *engine.E) error {
	println("initializing state")
	// Init Glow
	if err := gl.Init(); err != nil {
		return err
	}

	// Enable debug output
	if glfw3.ExtensionSupported("GL_ARB_debug_output") {
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS_ARB)
		gl.DebugMessageCallbackARB(gl.DebugProc(glDebugCallback), nil)
	}

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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
	//gl.Enable(gl.DITHER)
	gl.Enable(gl.MULTISAMPLE)
	gl.Enable(gl.POLYGON_SMOOTH)
	gl.Enable(gl.PROGRAM_POINT_SIZE)

	// Init VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	s.vertexArray = vao
	eng.Defer(func(){
		gl.DeleteVertexArrays(1, &vao)
	})

	// Set projection
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Initialize vertex buffer
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	s.vertexBuffer = buffer
	eng.Defer(func(){
		gl.DeleteBuffers(1, &buffer)
	})

	positionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Initialize space
	space := bnticles.New()
	for i := 0; i < PARTICLE_COUNT_DEFAULT; i++ {
		particle := space.MkParticle(PARTICLE_MASS_DEFAULT)
		particle.Position = randVect()
		particle.Velocity = randVect().Div(4)
	}
	s.space = space

	// Deinit is not part of interface, but can be deferred manually:
	eng.Defer(s.Deinit)

	return nil
}

func (s *MainState) Deinit() {
}

func (s *MainState) Update(dt float64) error {
	s.space.Step(dt)
	return nil
}

func (s *MainState) Render() error {
	program := s.program
	gl.UseProgram(program)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	gl.Flush()
	buffer := s.vertexBuffer
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		int(PARTICLE_COUNT_DEFAULT*unsafe.Sizeof([2]float32{})),
		nil,
		gl.STREAM_DRAW)
	dataPointer := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)
	if dataPointer == nil {
		return fmt.Errorf("buffer data pointer is nil")
	}
	data := (*[PARTICLE_COUNT_DEFAULT][2]float32)(dataPointer)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	for i := range data {
		p := s.space.Particle(i)
		if p == nil {
			return fmt.Errorf("no particle of id %d", i)
		}
		data[i][0] = float32(p.Position.X)
		data[i][1] = float32(p.Position.Y)
	}
	gl.UnmapBuffer(gl.ARRAY_BUFFER)

	gl.DrawArrays(gl.POINTS, 0, PARTICLE_COUNT_DEFAULT)

	return nil
}

func (s *MainState) KeyPress(key glfw3.Key) {
	return
}

func (s *MainState) KeyRelease(key glfw3.Key) {
	return
}

func (s *MainState) Resize(width int, height int) {
	println("resize!", width, height)
	aspect := float32(width) / float32(height)
	if width > height {
		projection = mgl32.Ortho2D(-aspect, aspect, -1, 1)
	} else {
		projection = mgl32.Ortho2D(-1, 1, -1/aspect, 1/aspect)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	return
}

func (s *MainState) Close() error {
	return nil
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
	{ // Profiling
		cpu, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer cpu.Close()
		pprof.StartCPUProfile(cpu)
		defer pprof.StopCPUProfile()
	}


	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	params := engine.Params{
		Version: [2]int{4, 3},
		Size:    [2]int{512, 512},
	}
	engine := engine.New()
	if err := engine.Initialize("TICKLES", params, &MainState{}); err != nil {
		panic(err)
	}
	if err := engine.Run(); err != nil {
		panic(err)
	}
	if err := engine.Shutdown(); err != nil {
		panic(err)
	}
}
