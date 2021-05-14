package objects

import "github.com/jevans40/psychic-spork/linmath"

//This is an interface that allows objects to set as their environment.
//The point is that the workers can handle all the updates for windows and input events.
//They can then relay the information to the objects at the objects request.
type ObjEnviroment interface {
	GetWindowSize() linmath.PSPoint
	GetEventCallback() EventCallback
	GetMousePosition() linmath.PSPoint
}
