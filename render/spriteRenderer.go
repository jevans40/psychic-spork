package render

type spriteRenderer struct {
	vertexBufferObject  uint32
	vertexArrayObject   uint32
	elementBufferObject uint32
	numOfSprites        int32
	programObject       uint32
	vert                []float32
	elem                []uint32
	allocation          []bool
}

func SpriteRendererFactory() Renderer {
	return &
}

func (thisRenderer *spriteRenderer) GetArraySprite(sprite Sprite) {

}

func (thisRenderer *spriteRenderer) allocate() {

}

func (thisRenderer *spriteRenderer) deallocate(num int32) {

}
