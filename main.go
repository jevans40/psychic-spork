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

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

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

	var vbo uint32

	// Generate 1 buffer, put the resulting identifier in vertexbuffer
	gl.GenBuffers(1, &vbo)
	// The following commands will talk about our 'vertexbuffer' buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Give our vertices to OpenGL
	gl.BufferData(gl.ARRAY_BUFFER, len(searialTriangle)*4, gl.Ptr(searialTriangle), gl.STATIC_DRAW)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	for !gameWindow.GetWindow().ShouldClose() {
		graphics.Update()

		render.Render(vbo, program)

		gameWindow.GetWindow().SwapBuffers()
		glfw.PollEvents()
	}

}
