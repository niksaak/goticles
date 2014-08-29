package engine

import "github.com/go-gl/glfw3"

// State initializes on creation and deinitializes on destruction.
type State interface {
	Init(*E) error
}

// Updater state updates its data with fixed timestep when active.
type Updater interface {
	Update(dt float64) error
}

// Renderer state renders itself on screen when active.
type Renderer interface {
	Render() error
}

// KeyPresser state handles keys which are pressed when it's active.
type KeyPresser interface {
	KeyPress(glfw3.Key)
}

// KeyReleaser state handles key releases when active.
type KeyReleaser interface {
	KeyRelease(glfw3.Key)
}

// Resizer state handles window resizes when active.
type Resizer interface {
	Resize(width int, height int)
}

// Closer state handles window closing.
type Closer interface {
	Close() error
}

// StateFull state is a state which satisfies all state interfaces.
type StateFull interface {
	State
	Updater
	Renderer
	KeyPresser
	KeyReleaser
	Resizer
	Closer
}
