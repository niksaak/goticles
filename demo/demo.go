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
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

func randFloat() float64 {
	return rand.Float64() - 0.5
}

func randVect() vect.V {
	return vect.V{randFloat(), randFloat()}
}

type Spacer interface {
	Particle(id int) *goticles.P
	MkParticle(mass float64) *goticles.P
	Step(dt float64)
}

const (
	PARTICLE_MASS_DEFAULT  = 100
	PARTICLE_COUNT_DEFAULT = 512
)

type MainState struct {
	space        Spacer
	vertexArray  uint32
	vertexBuffer uint32
	program      uint32
}

func (s *MainState) Init() error {
	println("initializing state")
	// Initialize rendering
	program, err := initGraphics()
	if err != nil {
		return err
	}
	s.program = program

	// Init VAO
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	s.vertexArray = vao

	// Set projection
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Initialize vertex buffer
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	s.vertexBuffer = buffer

	positionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(positionAttrib)
	gl.VertexAttribPointer(positionAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// Initialize space
	space := bnticles.New()
	for i := 0; i < PARTICLE_COUNT_DEFAULT; i++ {
		particle := space.MkParticle(PARTICLE_MASS_DEFAULT)
		particle.Position = randVect().Mul(0.5)
	}
	s.space = space

	return nil
}

func (s *MainState) Deinit() error {
	return nil
}

const accuracy = 1

func (s *MainState) Update(dt float64) error {
	s.space.Step(dt / accuracy)
	return nil
}

func unsafeEnslice(p unsafe.Pointer, size int, length int) unsafe.Pointer {
	sliceHeader := reflect.SliceHeader{uintptr(p), size * length, size * length}
	return unsafe.Pointer(&sliceHeader)
}

func (s *MainState) Render() error {
	program := s.program
	gl.UseProgram(program)

	gl.Clear(gl.COLOR_BUFFER_BIT)

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	buffer := s.vertexBuffer
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
	gl.InvalidateBufferData(buffer)
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

	for i := range data {
		p := s.space.Particle(i)
		if p == nil {
			return fmt.Errorf("no particle of id %d", i)
		}
		data[i][0] = float32(p.Position.X)
		data[i][1] = float32(p.Position.Y)
	}
	gl.UnmapBuffer(gl.ARRAY_BUFFER)

	gl.DrawArrays(gl.POINTS, 0, PARTICLE_COUNT_DEFAULT*2)

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

var projection = mgl32.Ident4()
var vertexSource = `
#version 430

uniform mat4 projection;

in vec2 position;

void main() {
	gl_Position.xy = position;
	gl_Position.z = 0;
	gl_Position.w = 1.0;
	gl_Position *= projection;
}
`
var fragmentSource = `
#version 430

in vec4 color;

void main() {
	gl_FragColor = vec4(0.3, 0.9, 0.3, 1.0);
}
`

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

func kindString(kind uint32) string {
	switch kind {
	case gl.VERTEX_SHADER:
		return "vertex"
	case gl.TESS_CONTROL_SHADER:
		return "tesselation control"
	case gl.TESS_EVALUATION_SHADER:
		return "tesselation evaluation"
	case gl.GEOMETRY_SHADER:
		return "geometry"
	case gl.FRAGMENT_SHADER:
		return "fragment"
	case gl.COMPUTE_SHADER:
		return "compute"
	default:
		return "unknown"
	}
}

func compileShader(source string, kind uint32) (uint32, error) {
	shader := gl.CreateShader(kind)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v shader:\n%s", kindString(kind), log)
	}

	return shader, nil
}

func newProgram(vertexSource, fragmentSource string) (uint32, error) {
	vertexSource += "\x00"
	fragmentSource += "\x00"

	vertexShader, err := compileShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(vertexShader)
	fragmentShader, err := compileShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	defer gl.DeleteShader(fragmentShader)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program:\n%s", log)
	}

	return program, nil
}

func initGraphics() (uint32, error) {
	// Init Glow
	if err := gl.Init(); err != nil {
		return 0, err
	}

	// Compile shaders
	program, err := newProgram(vertexSource, fragmentSource)
	if err != nil {
		return program, err
	}
	gl.UseProgram(program)

	// Enable debug output
	if glfw3.ExtensionSupported("GL_ARB_debug_output") {
		gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS_ARB)
		gl.DebugMessageCallbackARB(gl.DebugProc(glDebugCallback), nil)
	}

	return program, nil
}

func main() {
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
