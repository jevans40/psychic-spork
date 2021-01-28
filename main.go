package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/psychic-spork/graphics"
	"github.com/jevans40/psychic-spork/logging"
	"github.com/jevans40/psychic-spork/render"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	//Setup Logging
	logging.Initalize()

	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		logging.Log.Panic(err)
	}
	defer glfw.Terminate()

	gameWindow := graphics.NewWindow(1080, 920)
	if err != nil {
		logging.Log.Panic(err)
	}

	glfw.DetachCurrentContext()
	(gameWindow.GetWindow()).MakeContextCurrent()

	if err := gl.Init(); err != nil {
		logging.Log.Panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	logging.Log.Notice("Starting the Psychic-Spork game engine!")
	logging.Log.Infof("OpenGL version: %s", version)

	var nrAttributes int32
	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	logging.Log.Infof("Max Attributes Supported: %v", nrAttributes)

	//Get maximum texture size
	var MaximumTextureSize int32
	gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &MaximumTextureSize)
	logging.Log.Infof("Maximum Texture Size: %v", MaximumTextureSize)

	var MaximumTextureUnits int32
	gl.GetIntegerv(gl.MAX_TEXTURE_IMAGE_UNITS, &MaximumTextureUnits)
	logging.Log.Infof("Maximum Texture Units: %v", MaximumTextureUnits)

	//Set Renderer
	thisRender := render.SpriteRendererFactory()

	//Test Sprites move to render_test later
	sprite := render.SimpleSpriteFactory(&thisRender)
	sprite.Move(10, 10, 50)
	sprite.Resize(100, 100)
	sprite.Recolor(50, 255, 50, 255)
	sprite.SetTexSize(1, 1)

	sprite2 := render.SimpleSpriteFactory(&thisRender)
	sprite2.Move(200, 200, 50)
	sprite2.Resize(150, 200)
	sprite2.Recolor(255, 127, 0, 255)
	sprite2.SetTexSize(1, 1)
	sprite2.SetMap(1)

	sprite3 := render.SimpleSpriteFactory(&thisRender)
	sprite3.Move(800, 200, 50)
	sprite3.Resize(150, 200)
	sprite3.Recolor(125, 127, 0, 255)
	sprite3.SetTexSize(1, 1)
	sprite3.SetMap(1)

	for !gameWindow.GetWindow().ShouldClose() {
		graphics.Update()
		x, y := gameWindow.GetSize()
		thisRender.Render(int32(x), int32(y))
		r, g, b, a := sprite2.GetColor()
		sprite2.Recolor((r+1)%255, (g+1)%255, (b+1)%255, a)
		gameWindow.GetWindow().SwapBuffers()
		glfw.PollEvents()
	}

}
