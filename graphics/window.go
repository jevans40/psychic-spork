package graphics

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// GoWindow the struct that holds glfw window information
type GoWindow struct {
	lWindow *glfw.Window
	lErr    error
	lSize   [2]int
}

// NewWindow Initalizes the information into a goWindow
// x int : Size in the x dimension
// y int : Size in the y dimension
func NewWindow(x int, y int) *GoWindow {
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(x, y, "Testing", nil, nil)
	thisWindow := GoWindow{}
	thisWindow.lWindow = win
	thisWindow.lErr = err
	thisWindow.lSize = [2]int{x, y}
	return &thisWindow
}

//GetError returns the lErr feild associated with GoWindow
func (thisWindow *GoWindow) GetError() error {
	return thisWindow.lErr
}

//GetWindow returns the lWindow feild associated with GoWindow
func (thisWindow *GoWindow) GetWindow() *glfw.Window {
	return thisWindow.lWindow
}

//GetSize returns the lSize feild associated with GoWindow
func (thisWindow *GoWindow) GetSize() (x int, y int) {
	x, y = thisWindow.lWindow.GetSize()
	thisWindow.lSize[0] = x
	thisWindow.lSize[1] = y
	return
}
