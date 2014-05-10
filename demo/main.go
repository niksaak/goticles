package main

import (
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw3"
	//"github.com/go-gl/glh"
	"github.com/go-gl/gltext"
	. "github.com/niksaak/goticles"
	"github.com/niksaak/goticles/vect"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const (
	PARTICLE_COUNT = 4096
	PARTICLE_MASS = 1.0/2048
)

func randFloat05() float64 {
	return rand.Float64() - 0.5
}

func main() {
	rand.Seed(time.Now().UnixNano())
	space := MkSpace()
	for i := 0; i < PARTICLE_COUNT; i++ {
		particle := space.MkParticle(PARTICLE_MASS)
		r := rand.Float64()
		th := 2 * math.Pi * rand.Float64()
		particle.SetPosition(vect.Angle(th).Mul(r))
		px := randFloat05() * PARTICLE_MASS
		py := randFloat05() * PARTICLE_MASS
		particle.Momentum = particle.Momentum.Add(vect.V{px, py})
		particle.Recalculate()
	}
	screen, err := MkScreen(512, 512, "GO TICKLES")
	if err != nil {
		panic(err)
	}
	SimulateOnScreenFixed(space, 1.0/60, screen)
}

func SimulateOnScreenFixed(s *Space, step float64, screen *Screen) {
	const (
		STEPPING_MULTIPLIER = 1.0/2
	)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	screen.window.MakeContextCurrent()

	points := make([]vect.V, len(s.Particles))

	tick := time.Tick(time.Duration(step * float64(time.Second)))
	debug := time.Tick(1 * time.Second)

	var fps int

	for !screen.window.ShouldClose() {
		select {
		case <-tick:
			s.Step(step * STEPPING_MULTIPLIER)
			for i := range points {
				points[i] = s.Particles[i].Position
			}
			gl.Clear(gl.COLOR_BUFFER_BIT)
			gl.VertexPointer(2, gl.DOUBLE, 0, points)
			gl.DrawArrays(gl.POINTS, 0, len(points))

			screen.window.SwapBuffers()
			glfw3.PollEvents()
			fps++
		case <-debug:
			log.Printf("fps: %d, %v", fps, s.Particle(0))
			fps = 0
		}
	}
}

type Screen struct {
	window *glfw3.Window
	font *gltext.Font
}

func MkScreen(width, height int, title string) (*Screen, error) {
	return new(Screen).Init(width, height, title)
}

func setupGFX(width, height int, title string) (*glfw3.Window, error) {
	// Lock OS thread while setting up.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Setup GLFW
	if err := glfw3.Init(); err != nil {
		return nil, err
	}
	glfw3.WindowHint(glfw3.Resizable, 1)
	glfw3.WindowHint(glfw3.ContextVersionMajor, 2)
	glfw3.WindowHint(glfw3.ContextVersionMinor, 1)
	window, err := glfw3.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	window.SetSizeCallback(func(window *glfw3.Window, width, height int) {
		log.Printf("fixing aspect ratio for %dx%d", width, height)
		a := float64(width) / float64(height)

		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Viewport(0, 0, width, height)

		if width > height {
			gl.Ortho(-a, a, -1, 1, -1, 1)
		} else {
			gl.Ortho(-1, 1, 1/-1, 1/a, -1, 1)
		}

		gl.MatrixMode(gl.MODELVIEW)
	})
	window.SetCloseCallback(func(window *glfw3.Window) {
		log.Println("closing window")
	})

	// Setup OpenGL
	gl.ClearColor(0, 0, 0, 1)
	gl.LineWidth(1.2)
	gl.PointSize(1)

	gl.Enable(gl.ALPHA_TEST)
	gl.AlphaFunc(gl.LEQUAL, 1)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.EnableClientState(gl.VERTEX_ARRAY)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.PushMatrix()

	return window, nil
}

func setupFont(scale int32) (*gltext.Font, error) {
	path, err := gas.Abs(
		"code.google.com/p/freetype-go/testdata/luxisr.ttf")
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return gltext.LoadTruetype(file, scale, 32, 127, gltext.LeftToRight)
}

func (s *Screen) Init(width, height int, title string) (*Screen, error) {
	const(
		DEFAULT_WIDTH = 512
		DEFAULT_HEIGHT = 512
	)
	if width <= 0 {
		width = DEFAULT_WIDTH
	}
	if height <= 0 {
		height = DEFAULT_HEIGHT
	}

	var err error

	// Set window for screen
	s.window, err = setupGFX(width, height, title)
	if err != nil {
		return nil, err
	}

	// Setup font
	s.font, err = setupFont(12)
	if err != nil {
		return nil, err
	}

	return s, nil
}

