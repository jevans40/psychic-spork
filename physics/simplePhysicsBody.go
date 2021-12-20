package physics

import ()

type SimplePhysicsBody interface {
	GetDeltas(deltaTime int) (x float32, y float32, z float32)
	SetVelocity(x float32, y float32, z float32)
	SetAcceleration(x float32, y float32, z float32)
	SetPosition(x float32, y float32, z float32)
	SetBounds(x float32, y float32, z float32)

	GetVelocity() (x float32, y float32, z float32)
	GetAcceleration() (x float32, y float32, z float32)
	GetPosition() (x float32, y float32, z float32)
	GetBounds() (x float32, y float32, z float32)
}
