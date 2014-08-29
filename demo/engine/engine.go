package engine

import (
	"fmt"
	"github.com/go-gl/glfw3"
	"time"
)

// E structure contains data necessary for engine operation.
type E struct {
	window *glfw3.Window
	state  State
	defers []func()
}

// New returns pointer to uninitialized new engine.
func New() *E {
	return new(E)
}

// Params struct contains engine initialization parameters.
type Params struct {
	Version [2]int // major, minor
	Size    [2]int // width, height
}

// Initialize creates new window 
func (e *E) Initialize(title string, params Params, state State) error {
	if !glfw3.Init() {
		return fmt.Errorf("unable to initialize glfw")
	}
	if e.window != nil {
		return fmt.Errorf("Initialize error: window is not nil")
	}

	glfw3.SetErrorCallback(OnError)

	glfw3.WindowHint(glfw3.Resizable, glfw3.True)
	glfw3.WindowHint(glfw3.ContextVersionMajor, params.Version[0])
	glfw3.WindowHint(glfw3.ContextVersionMinor, params.Version[1])
	glfw3.WindowHint(glfw3.OpenglDebugContext, glfw3.True)

	size := params.Size
	window, err := glfw3.CreateWindow(size[0], size[1], title, nil, nil)
	if err != nil {
		return err
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(e.OnKey)
	window.SetFramebufferSizeCallback(e.OnFramebufferResize)

	e.window = window
	e.state = state

	return nil
}

// IsInitialized returns true if engine is fully initialized.
func (e *E) IsInitialized() bool {
	return e.window != nil && e.state != nil
}

// Run starts the main loop.
func (e *E) Run() error {
	if !e.IsInitialized() {
		return fmt.Errorf("not initialized")
	}

	if err := e.state.Init(e); err != nil {
		return err
	}

	if state, ok := e.state.(Resizer); ok {
		width, height := e.window.GetFramebufferSize()
		state.Resize(width, height)
	}

	ticker := time.Tick(1 * time.Second / 60)
	for _ = range ticker {
		glfw3.PollEvents()

		if e.window.ShouldClose() {
			if state, ok := e.state.(Closer); ok {
				if err := state.Close(); err != nil {
					return err
				}
			}
			return nil
		}

		if state, ok := e.state.(Updater); ok {
			if err := state.Update(1.0 / 60); err != nil {
				return err
			}
		}
		if state, ok := e.state.(Renderer); ok {
			if err := state.Render(); err != nil {
				return err
			}
		}

		e.window.SwapBuffers()
	}

	return nil
}

// Shutdown deinitializes all states and closes all windows.
func (e *E) Shutdown() (err error) {
	if !e.IsInitialized() {
		return NotInitializedError("shutdown")
	}
	e.window.SetShouldClose(true)
	err = e.executeDefers()
	glfw3.Terminate()
	*e = E{}
	return err
}

func (e *E) OnKey(
	w *glfw3.Window,
	key glfw3.Key,
	scancode int,
	action glfw3.Action,
	mods glfw3.ModifierKey) {

	switch action {
	case glfw3.Press, glfw3.Repeat:
		if state, ok := e.state.(KeyPresser); ok {
			state.KeyPress(key)
		}
	case glfw3.Release:
		if state, ok := e.state.(KeyReleaser); ok {
			state.KeyRelease(key)
		}
	}
}

func (e *E) OnFramebufferResize(w *glfw3.Window, width int, height int) {
	if state, ok := e.state.(Resizer); ok {
		state.Resize(width, height)
	}
}

func OnError(code glfw3.ErrorCode, desc string) {
	panic(desc)
}

type NotInitializedError string

func (e NotInitializedError) Error() string {
	s := "not initialized"
	if e != "" {
		s += ": " + string(e)
	}
	return s
}
