package objects

import "github.com/jevans40/psychic-spork/event"

type EventCallback func(e event.UpdateEvent)

type Object interface {
	SetEnviroment(env ObjEnviroment)
	SendEvent(event.UpdateEvent)
	Update(int)
	Render([]float32)
}
