package objects

import (
	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/linmath"
)

//This is the bare minimum object
//It just allows you to update it and render the vertex array
//Also allows you to interact with the UpdateEvent System
var _ Object = (*BaseObject)(nil)

type BaseObject struct {
	env             ObjEnviroment
	size            [2]float32
	position        [3]float32
	texturePosition [2]float32
	textureSize     [2]float32
	textureMap      uint32
	color           [4]uint8
	verts           linmath.FourVertice
	floats          []float32
}

func (b *BaseObject) Update(int) {
	//Eat Update
}

func (b *BaseObject) Init() {
	if b.verts == nil {

		b.verts = linmath.EmptyFourVertice()

	}
}

func (b *BaseObject) Render(toreturn []float32) {
	if b.verts == nil {

		b.verts = linmath.EmptyFourVertice()
	}
	b.verts.SetTexX(b.texturePosition[0], b.textureSize[0])
	b.verts.SetTexY(b.texturePosition[1], b.textureSize[1])
	b.verts.SetX(b.position[0], b.size[0])
	b.verts.SetY(b.position[1], b.size[1])
	b.verts.SetZ(b.position[2])

	copy(toreturn, b.verts.ToFloats()[:])
}

func (b *BaseObject) SetEnviroment(env ObjEnviroment) {
	b.env = env
	//Set Callback
}

func (b *BaseObject) SendEvent(e event.UpdateEvent) {
	//Eat Event
}

func (b *BaseObject) Resize(x float32, y float32) {
	b.size = [2]float32{x, y}
}

func (b *BaseObject) Move(x float32, y float32, z float32) {
	b.position = [3]float32{x, y, z}
}

func (b *BaseObject) Recolor(r uint8, g uint8, bl uint8, a uint8) {
	b.color = [4]uint8{a, bl, g, r}
	b.verts.SetColor(b.color)
}

func (b *BaseObject) SetTexPos(x float32, y float32) {
	b.texturePosition = [2]float32{x, y}
}

func (b *BaseObject) SetTexSize(x float32, y float32) {
	b.textureSize = [2]float32{x, y}
}

func (b *BaseObject) SetMap(newMap uint32) {
	b.textureMap = newMap
	b.verts.SetMap(b.textureMap)
}

func (b *BaseObject) GetSize() (x, y float32) {
	x, y = b.size[0], b.size[1]
	return
}

func (b *BaseObject) GetPos() (x, y, z float32) {
	x, y, z = b.position[0], b.position[1], b.position[2]
	return
}

func (b *BaseObject) GetColor() (r, g, bl, a uint8) {
	r, g, bl, a = b.color[0], b.color[1], b.color[2], b.color[3]
	return
}

func (b *BaseObject) GetTexPos() (x, y float32) {
	x, y = b.texturePosition[0], b.texturePosition[1]
	return
}

func (b *BaseObject) GetTexSize() (x, y float32) {
	x, y = b.textureSize[0], b.textureSize[1]
	return
}

func (b *BaseObject) GetMap() (texMap uint32) {
	texMap = b.textureMap
	return
}

func (b *BaseObject) MakeObject(object Object) {
	newObjectEvent := event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
		Receiver: -1,
		Sender:   -1,
		Event:    event.UpdateEvent_NewObjectEvent{object}}
	b.env.GetEventCallback()(newObjectEvent)
}
