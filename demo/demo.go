package main

import (
	"math/rand"
	"github.com/niksaak/goticles/vect"
	"github.com/niksaak/goticles/demo/engine"
	"github.com/go-gl/glfw3"
	"time"
)

func randFloat() float64 {
	return rand.Float64() - 0.5
}

func randVect() vect.V {
	return vect.V{randFloat(), randFloat()}
}

type Stater interface{
	Init() error
	Deinit() error
	Update(dt float64) error
	Draw() error
}

type MainState struct {
}

func (s *MainState) Init() error {
	return nil
}

func (s *MainState) Deinit() error {
	return nil
}

func (s *MainState) Update(dt time.Duration) error {
	return nil
}

func (s *MainState) Render() error {
	return nil
}

func (s *MainState) KeyPress(key glfw3.Key) {
	return
}

func (s *MainState) KeyRelease(key glfw3.Key) {
	return
}

func (s *MainState) Resize(width int, height int) {
	return
}

func (s *MainState) Close() error {
	return nil
}

var _ engine.StateFull = &MainState{}

func main() {
	params := engine.Params{
		Version: [2]int{4,3},
		Size: [2]int{512,512},
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

