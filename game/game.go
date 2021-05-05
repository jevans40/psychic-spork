package game

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/psychic-spork/graphics"
	"github.com/jevans40/psychic-spork/logging"
	"github.com/jevans40/psychic-spork/render"
)

//Handles setting up all the needed components and creates a clock to run the game.
//This is the entrypoint for the rest of the program
type Game struct {
	window   *graphics.GoWindow
	renderer render.SpriteRenderer
}

func (g *Game) Init() error {
	//Initalize Logging for the game engine
	logging.Initalize()

	//Lock the main thread since it is required for GLFW to function properly
	runtime.LockOSThread()

	//Initialize glfw
	err := glfw.Init()
	if err != nil {
		return err
	}
	return nil
}

//Starts the game
//Warning, call this function only once you are ready for the game to start.
//This is an infinite loop and should be the last line in the program
func (g *Game) Start() {
	fmt.Print("Starting")

	//Generate the window for the game
	window, err := graphics.NewWindow(1080, 920)
	g.window = window
	if err != nil {
		logging.Log.Panic(err)
	}
	defer glfw.Terminate()
	glfw.WaitEvents()

	glfw.DetachCurrentContext()
	(g.window.GetWindow()).MakeContextCurrent()

	if err := gl.Init(); err != nil {
		logging.Log.Panic(err)
	}

	//Log Debug Info to the console
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

	g.renderer = render.SpriteRendererFactory()
}
