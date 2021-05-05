package game

import (
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/graphics"
	"github.com/jevans40/psychic-spork/logging"
	"github.com/jevans40/psychic-spork/render"
	"github.com/jevans40/psychic-spork/update"
)

const Version = "0.1"

//Handles setting up all the needed components and creates a clock to run the game.
//This is the entrypoint for the rest of the program
type Game struct {
	window               *graphics.GoWindow
	renderer             render.SpriteRenderer
	coordinator          *update.Coordinator
	EventChannel         chan []event.UpdateEvent
	RenderChannel        chan []float32
	CommunicationChannel chan int
}

func (g *Game) Init() error {
	//Initalize Channels for component communication
	g.EventChannel = make(chan []event.UpdateEvent)
	g.RenderChannel = make(chan []float32)
	g.CommunicationChannel = make(chan int)

	//Create the renderer
	g.renderer = render.SpriteRendererFactory()

	//Create the coordinator
	g.coordinator = update.CoordinatorFactory(g.EventChannel, g.RenderChannel)

	return nil
}

//Starts the game
//Warning, call this function only once you are ready for the game to start.
//This is an infinite loop and should be the last line in the program
func (g *Game) Start() {

	//Initalize Logging for the game engine
	logging.Initalize()

	//Lock the main thread since it is required for GLFW to function properly
	runtime.LockOSThread()

	//Initialize glfw
	err := glfw.Init()
	if err != nil {
		logging.Log.Panic(err)
	}

	//Create the window
	window, err := graphics.NewWindow(1080, 920)
	g.window = window
	if err != nil {
		logging.Log.Panic(err)
	}
	defer glfw.Terminate()
	glfw.WaitEvents()

	//Initialize window context
	glfw.DetachCurrentContext()
	(g.window.GetWindow()).MakeContextCurrent()

	//Initialize open-gl
	if err := gl.Init(); err != nil {
		logging.Log.Panic(err)
	}

	fmt.Print("Starting")

	//Generate the window for the game

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

	//Log game version
	logging.Log.Infof("Psychic Spork Version: %s", Version)

	syncClock := Clock{}
	syncClock.SetChannels(g.CommunicationChannel)
	syncClock.SetDurations(time.Millisecond*time.Duration(16), time.Millisecond*time.Duration(5))

	go syncClock.Start()
	go g.Render()
	go g.coordinator.Start(g.CommunicationChannel)

	event.EventLoop()
}

func (g *Game) Render() {
	//Old vs new buffer
	var Buffer []float32
	numObjects := 0
	var sprites []*render.VertexRenderable
	for {
		Buffer = <-g.RenderChannel

		numObjects = len(Buffer) / 28
		for len(sprites) < numObjects {
			sprites = append(sprites, render.VertexSpriteFactory(&g.renderer))
		}
		for len(sprites) > numObjects {
			sprites = sprites[0 : len(sprites)-2]
		}
		for i, v := range sprites {
			v.SetVerticies(Buffer[i*28 : (i+1)*28])
		}
		x, y := g.window.GetSize()
		g.renderer.Render(int32(x), int32(y))
		g.window.GetWindow().SwapBuffers()
	}
}
