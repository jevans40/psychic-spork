package objects

import (
	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/linmath"
)

//This is the bare minimum object
//It just allows you to update it and render the vertex array
var _ Object = (*BaseObject)(nil)

type BaseObject struct {
	eventCallback   EventCallback
	size            [2]float32
	position        [3]float32
	texturePosition [2]float32
	textureSize     [2]float32
	textureMap      uint32
	color           [4]uint8
}

func (b *BaseObject) Update(int) {
	//Eat Update
}

func (b *BaseObject) Render() []float32 {
	var verticeData [28]float32
	vert := make([]linmath.Vertice, 4)
	for i, _ := range vert {
		vert[i] = linmath.EmptyVertice()
	}

	for i := 0; i < 4; i++ {
		vert[i].SetColor(b.color)
		vert[i].SetMap(b.textureMap)
		vert[i].SetTexX(b.texturePosition[0] + b.textureSize[0]*float32((i)%2))
		vert[i].SetTexY(b.texturePosition[1] + b.textureSize[1]*float32(int32((i)/2)%2))
		vert[i].SetX(b.position[0] + b.size[0]*float32((i)%2))
		vert[i].SetY(b.position[1] + b.size[1]*float32(int32((i)/2)%2))
		vert[i].SetZ(b.position[2])
	}

	for i, v := range vert {
		floats := v.ToFloats()
		copy(verticeData[7*i:7*i+7], floats[:])
	}
	return verticeData[:]
}

func (b *BaseObject) SetEventCallback(call EventCallback) {
	b.eventCallback = call
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
	b.color = [4]uint8{r, g, bl, a}
}

func (b *BaseObject) SetTexPos(x float32, y float32) {
	b.texturePosition = [2]float32{x, y}
}

func (b *BaseObject) SetTexSize(x float32, y float32) {
	b.textureSize = [2]float32{x, y}
}

func (b *BaseObject) SetMap(newMap uint32) {
	b.textureMap = newMap
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
