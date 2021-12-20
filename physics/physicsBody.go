package physics

type PhysicsBody interface {
	GetDeltas(deltaTime int) (x int, y int, z int)
}
