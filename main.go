package main

import (
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

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

var cpuprofile = flag.String("cpuprofile", "cputest", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "memtest", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	//Setup Logging
	logging.Initalize()

	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		logging.Log.Panic(err)
	}
	defer glfw.Terminate()

	gameWindow, err := graphics.NewWindow(1080, 920)
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
	sprite3.Move(800, 200, 49)
	sprite3.Resize(150, 200)
	sprite3.Recolor(125, 127, 0, 255)
	sprite3.SetTexSize(1, 1)
	sprite3.SetMap(1)

	Spritenum := 0
	sprites := []render.Sprite{}
	count := 1
	lastframe := time.Now()
	lastUpdate := time.Now()
	lastSecond := time.Now()
	fps := 0
	for !gameWindow.GetWindow().ShouldClose() {
		x, y := gameWindow.GetSize()
		graphics.Update()
		if time.Duration(5)*time.Millisecond < time.Since(lastUpdate) {
			lastUpdate = time.Now()
			Spritenum++
			if Spritenum%1 == 0 {
				count++
				spritetmp := render.SimpleSpriteFactory(&thisRender)
				spritetmp.Move(22, 33, 50/float32(count)+10)
				spritetmp.Resize(float32(int(Spritenum*1)%100), 100)
				spritetmp.Recolor(uint8(Spritenum*17)%255, uint8(Spritenum*1)%255, uint8(Spritenum*1)%255, uint8(Spritenum*4)%255)
				spritetmp.SetTexSize(1, 1)
				sprites = append(sprites, spritetmp)
			}

			if Spritenum%3 == 0 {
				count++
				spritetmp := render.SimpleSpriteFactory(&thisRender)
				spritetmp.Move(15, 10, 50/float32(count)+10)
				spritetmp.Resize(float32(int(Spritenum*57)%100), float32(int(Spritenum*2)%100))
				spritetmp.Recolor(uint8(Spritenum*23)%255, uint8(Spritenum*2)%255, uint8(Spritenum*7)%255, uint8(Spritenum*4)%255)
				spritetmp.SetTexSize(1, 1)
				sprites = append(sprites, spritetmp)
			}

			if Spritenum%5 == 0 {
				count++
				spritetmp := render.SimpleSpriteFactory(&thisRender)
				spritetmp.Move(10, 10, 50/float32(count)+10)
				spritetmp.Resize(float32(int(Spritenum*37)%100), float32(int(Spritenum*4)%100))
				spritetmp.Recolor(uint8(Spritenum*19)%255, uint8(Spritenum*4)%255, uint8(Spritenum*3)%255, uint8(Spritenum*1)%255)
				spritetmp.SetTexSize(1, 1)
				sprites = append(sprites, spritetmp)
			}

			if Spritenum%9 == 0 {
				count++
				spritetmp := render.SimpleSpriteFactory(&thisRender)
				spritetmp.Move(10, 10, 50/float32(count)+10)
				spritetmp.Resize(float32(int(Spritenum*13)%100), float32(int(Spritenum*8)%100))
				spritetmp.Recolor(uint8(Spritenum*29)%255, uint8(Spritenum*5)%255, uint8(Spritenum*10)%255, uint8(Spritenum*11)%255)
				spritetmp.SetTexSize(1, 1)
				sprites = append(sprites, spritetmp)
			}

			if Spritenum%11 == 0 {
				count++
				spritetmp := render.SimpleSpriteFactory(&thisRender)
				spritetmp.Move(10, 10, 50/float32(count))
				spritetmp.Resize(float32(int(Spritenum*5)%100), float32(int(Spritenum*16)%100))
				spritetmp.Recolor(uint8(Spritenum*31)%255, uint8(Spritenum*1)%255, uint8(Spritenum*13)%255, uint8(Spritenum*9)%255)
				spritetmp.SetTexSize(1, 1)
				sprites = append(sprites, spritetmp)
			}
			for _, v := range sprites {
				fx, fy, fz := v.GetPos()
				v.Move(float32((int(fx)+5)%x), float32((int(fy)+11)%y), fz)
			}

			r, g, b, a := sprite2.GetColor()
			sprite2.Recolor((r+1)%255, (g+1)%255, (b+1)%255, a)
			sx, sy, sz := sprite.GetPos()
			sprite.Move(float32((int(sx)+5)%x), float32((int(sy)+11)%y), sz)
			s1x, s1y, s1z := sprite2.GetPos()
			sprite2.Move(float32((int(s1x)+10)%x), float32((int(s1y)+7)%y), s1z)
		}
		if time.Duration(8)*time.Millisecond < time.Since(lastframe) {
			lastframe = time.Now()
			fps++
			thisRender.Render(int32(x), int32(y))
			gameWindow.GetWindow().SwapBuffers()
			glfw.PollEvents()
		}
		if time.Duration(1)*time.Second < time.Since(lastSecond) {
			lastSecond = time.Now()
			logging.Log.Noticef("FPS:%v, Sprites: %v", fps, count)
			fps = 0
		}
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
