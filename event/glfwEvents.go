package event

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

//This package gets mouse/keyboard input from the mouse
//Later support for glfw handled mouse and joystick
//Other events occur here aswell.

//Game Event is the default handler for events that happen
type GameEvent struct {
}

type Key struct {
}

type Mouse struct {
}

func EventLoop() {
	for {
		glfw.WaitEvents()
	}
}
