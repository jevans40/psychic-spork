package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/psychic-spork/linmath"
	"github.com/jevans40/psychic-spork/linmath"
)

type SpriteRenderer struct {
	vertexBufferObject  uint32
	vertexArrayObject   uint32
	elementBufferObject uint32
	numOfSprites        int32
	programObject       uint32
	vert                []float32

	//The order of elements should be 0,1,2,1,2,3
	elem []uint32

	allocation []bool
}

func SpriteRendererFactory( /*SourceAtlas ImageAtlas*/ ) SpriteRenderer {
	//size := SourceAtlas.imageSize
func SpriteRendererFactory() SpriteRenderer {

	//First Generate a Vertex Array to bind the vertex buffer object and the element buffer object to
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//Generate the vertex buffer object
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	//Generate the element buffer object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(0)
	//Unbind the vertex array

	program, err := CreateDefaultProgram()
	if err != nil {
		panic(err)
	}
	var vert []float32
	var elem []uint32
	var allocation []bool

	//Make the sprite atlas

	newRenderer := SpriteRenderer{vbo, vao, ebo, 0, program, vert, elem, allocation}
	newRenderer.init()
	return newRenderer
}

func (thisRenderer *SpriteRenderer) Render(width, height int32) {
	// 1st attribute buffer : vertices
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0, 0, width, height)

	orthomat := linmath.NewOrthoMat4f(float32(height), 0, 0, float32(width), 1, 0)
	gl.UseProgram(thisRenderer.programObject)
	loc := gl.GetUniformLocation(thisRenderer.programObject, gl.Str("MVT"+"\x00"))
	mat := orthomat.ToFloats()
	gl.UniformMatrix4fv(loc, 1, false, &mat[0])
	thisRenderer.bind()
	gl.DrawElements(gl.TRIANGLES, 6*thisRenderer.numOfSprites, gl.UNSIGNED_INT, gl.PtrOffset(0)) // Starting from vertex 0; 3 vertices total -> 1 triangle
	thisRenderer.unbind()

}

func (thisRenderer *SpriteRenderer) init() {
	thisRenderer.allocation = append(thisRenderer.allocation, false)
	thisRenderer.vert = append(thisRenderer.vert, make([]float32, 28)...)
	thisRenderer.elem = append(thisRenderer.elem, make([]uint32, 6)...)

	thisRenderer.bind()
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.2, 0.1, 0.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 28, gl.PtrOffset(0))

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 28, gl.PtrOffset(12))

	gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, 28, gl.PtrOffset(20))

	gl.VertexAttribPointer(3, 1, gl.INT, false, 28, gl.PtrOffset(24))

	thisRenderer.unbind()
}

func (thisRenderer *SpriteRenderer) bind() {

	gl.BindVertexArray(thisRenderer.vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, thisRenderer.vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(thisRenderer.vert)*4, gl.Ptr(thisRenderer.vert), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, thisRenderer.elementBufferObject)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(thisRenderer.elem)*4, gl.Ptr(thisRenderer.elem), gl.STATIC_DRAW)

}

func (thisRenderer *SpriteRenderer) unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (thisRenderer *SpriteRenderer) GetArraySprite() (returnSlice []float32, spriteNum int32) {
	slot := uint32(thisRenderer.allocate())
	spriteNum = int32(slot)
	returnSlice = thisRenderer.vert[slot*28 : (slot+1)*28]
	newElem := []uint32{slot * 4, slot*4 + 1, slot*4 + 2, slot*4 + 2, slot*4 + 3, slot*4 + 1}
	thisRenderer.numOfSprites = thisRenderer.numOfSprites + 1
	copy(thisRenderer.elem[slot*6:(slot+1)*6], newElem)
	return
}

func (thisRenderer *SpriteRenderer) allocate() int32 {
	for i, v := range thisRenderer.allocation {
		if !v {
			thisRenderer.allocation[i] = true
			return int32(i)
		}
	}

	thisRenderer.allocation = append(thisRenderer.allocation, true)
	thisRenderer.vert = append(thisRenderer.vert, make([]float32, 28)...)
	thisRenderer.elem = append(thisRenderer.elem, make([]uint32, 6)...)
	return int32(len(thisRenderer.allocation) - 1)
}

func (thisRenderer *SpriteRenderer) DeallocateSprite(num int32) {
	newVert := make([]float32, 28)
	copy(thisRenderer.vert[num*28:(num+1)*28], newVert)
	newElem := make([]uint32, 6)
	copy(thisRenderer.elem[num*6:(num+1)*6], newElem)

	thisRenderer.allocation[num] = false
}
