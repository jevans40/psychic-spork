package graphics

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

//Renderer is the base renderer object
type Renderer struct {
	vertexBufferObject uint32
	numOfTriangles     int
	programObject      uint32
	vert               []float32
}

//Render Main render loop of the game
func Render(vbo, program uint32) {
	// 1st attribute buffer : vertices
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.EnableVertexAttribArray(0)
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(
		0,               // attribute 0. No particular reason for 0, but must match the layout in the shader.
		3,               // size
		gl.FLOAT,        // type
		false,           // normalized?
		0,               // stride
		gl.PtrOffset(0), // array buffer offset
	)
	// Draw the triangle !
	gl.DrawArrays(gl.TRIANGLES, 0, 3) // Starting from vertex 0; 3 vertices total -> 1 triangle
	gl.DisableVertexAttribArray(0)
}

//func (thisRenderer *Renderer) AddVertex()
