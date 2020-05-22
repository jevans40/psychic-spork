package main

import (
	"fmt"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/graphics"
	"github.com/jevans40/linmath"
	"github.com/jevans40/render"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	gameWindow := graphics.NewWindow(600, 400)
	if err != nil {
		panic(err)
	}

	glfw.DetachCurrentContext()
	(gameWindow.GetWindow()).MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("Starting the Psychic-Spork game engine!")
	fmt.Println("OpenGL version", version)

	triangle := make([]linmath.Vector3, 3)

	triangle[0] = linmath.NewVector3(-1, -1, 0)
	triangle[1] = linmath.NewVector3(1, -1, 0)
	triangle[2] = linmath.NewVector3(0, 1, 0)

	searialTriangle := []float32{}
	for _, v := range triangle {
		x := v.ToFloats()
		searialTriangle = append(searialTriangle, x[:]...)
	}

	program, err := render.CreateDefaultProgram()
	if err != nil {
		panic(err)
	}

	thisRender := render.NewTestRenderer(program)

	for !gameWindow.GetWindow().ShouldClose() {
		graphics.Update()
		x, y := gameWindow.GetSize()
		thisRender.Render(int32(x), int32(y))

		gameWindow.GetWindow().SwapBuffers()
		glfw.PollEvents()
	}

}
