package main

import (
	"fmt"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/psychic-spork/graphics"
	"github.com/psychic-spork/linmath"
	"github.com/psychic-spork/render"
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

	gameWindow := graphics.NewWindow(1080, 920)
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

	thisRender := render.SpriteRendererFactory()

	sprite := render.SimpleSpriteFactory(&thisRender)
	sprite.Move(10, 10, 50)
	sprite.Resize(100, 100)
	sprite.Recolor(0, 255, 0, 255)

	sprite2 := render.SimpleSpriteFactory(&thisRender)
	sprite2.Move(200, 200, 50)
	sprite2.Resize(150, 200)
	sprite2.Recolor(255, 255, 0, 255)

	for !gameWindow.GetWindow().ShouldClose() {
		graphics.Update()
		x, y := gameWindow.GetSize()
		thisRender.Render(int32(x), int32(y))

		gameWindow.GetWindow().SwapBuffers()
		glfw.PollEvents()
	}

}
