package game

import (
	"fmt"
	"runtime"
	"runtime/debug"
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
	coordinator          *update.Coordinator
	EventChannel         chan []event.UpdateEvent
	RenderChannel        chan []float32
	CommunicationChannel chan int
}

func (g *Game) Init() error {
	//Initalize Channels for component communication
	g.EventChannel = make(chan []event.UpdateEvent, runtime.NumCPU())
	g.RenderChannel = make(chan []float32, runtime.NumCPU())
	g.CommunicationChannel = make(chan int, runtime.NumCPU())

	//We have two forever blocking locked threads, so its okay to spawn a few more that the scheduler can play with.
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	debug.SetGCPercent(200)

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

	event.EventsInit()

	//g.window.GetWindow().SetKeyCallback(glfw.KeyCallback(event.EventKeyCallback))

	//Create the coordinator
	g.coordinator = update.CoordinatorFactory(g.EventChannel, g.RenderChannel)

	syncClock := Clock{}
	syncClock.SetChannels(g.CommunicationChannel)
	syncClock.SetDurations(time.Millisecond*time.Duration(16), time.Millisecond*time.Duration(16))

	go syncClock.Start()
	go g.Render()

	eventloop := make(chan event.UpdateEvent)
	go g.coordinator.Start(g.CommunicationChannel, eventloop)
	go event.EventSubscriberLoop(eventloop)

	event.EventLoop()
}

func (g *Game) Render() {
	//Old vs new buffer

	//THIS HAS TO BE IN THE SAME THREAD AS THE OTHER RENDERING
	//Initialize open-gl
	runtime.LockOSThread()
	glfw.DetachCurrentContext()
	(g.window.GetWindow()).MakeContextCurrent()
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
	var Buffer []float32
	numObjects := 0
	var sprites []*render.VertexRenderable
	renderer := render.SpriteRendererFactory()
	for {
		//Create the renderer
		Buffer = <-g.RenderChannel

		numObjects = len(Buffer) / 28
		for len(sprites) < numObjects {
			sprites = append(sprites, render.VertexSpriteFactory(&renderer))
		}
		for len(sprites) > numObjects {
			sprites = sprites[0 : len(sprites)-2]
		}
		for i, v := range sprites {
			v.SetVerticies(Buffer[i*28 : (i+1)*28])
		}
		x, y := g.window.GetSize()
		renderer.Render(int32(x), int32(y))
		g.window.GetWindow().SwapBuffers()
	}
}
