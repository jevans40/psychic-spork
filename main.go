package main

import (
	"fmt"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/graphics"
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

	for !gameWindow.GetWindow().ShouldClose() {
		//GL Goes here
		gameWindow.GetWindow().SwapBuffers()
		glfw.PollEvents()
	}

}
