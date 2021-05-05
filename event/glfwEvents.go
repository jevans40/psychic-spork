package event

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

//This package gets mouse/keyboard input from the mouse
//Later support for glfw handled mouse and joystick
//Other events occur here aswell.

//Game Event is the default handler for events that happen
type GameEvent struct {
}

type Key struct {
	KeyName     string
	KeyScancode int
	KeyCode     int
}

type Mouse struct {
}

//A map of keynames to key structs
var keyStore map[string]Key

//A map of scancodes to keystructs
var keyInverse map[int]Key

var keyListeners map[string](map[int]chan int)
var keyListenerLock sync.Mutex

//All events are lossy, events that cant be sent will just be ignored.
func EventKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	keyListenerLock.Lock()
	for _, v := range keyListeners[keyInverse[scancode].KeyName] {
		select {
		case v <- int(action):
			//If we cant send a message then don't. Not our problem
		default:
		}
	}
	keyListenerLock.Unlock()
}

func EventSubscriberLoop(e chan UpdateEvent) {
	for {
		event := <-e
		if event.EventCode == UpdateEvent_Subscribe {
			ev := (event.Event).(UpdateEvent_SubscribeEvent)
			keyListenerLock.Lock()
			keyListeners[ev.EventSubscriptionName][event.Sender] = ev.ListeningChannel
			keyListenerLock.Unlock()
		} else if event.EventCode == UpdateEvent_UnSubscribe {
			ev := (event.Event).(UpdateEvent_UnSubscribeEvent)
			keyListenerLock.Lock()
			close(keyListeners[ev.EventSubscriptionName][event.Sender])
			delete(keyListeners[ev.EventSubscriptionName], event.Sender)
			keyListenerLock.Unlock()
		}

	}
}

func EventLoop() {
	for {
		glfw.WaitEvents()
	}
}

//Must be called before the main loop but after glfw is initialized
func EventsInit() {
	keyStore = make(map[string]Key)
	//You know.. there is probably a smart way to do this....
	//But.....
	keyStore["Unknown"] = Key{"Unknown", glfw.GetKeyScancode(glfw.KeyUnknown), int(glfw.KeyUnknown)}
	keyStore["Space"] = Key{"Space", glfw.GetKeyScancode(glfw.KeySpace), int(glfw.KeySpace)}
	/*
		KeyApostrophe   Key = C.GLFW_KEY_APOSTROPHE
		KeyComma        Key = C.GLFW_KEY_COMMA
		KeyMinus        Key = C.GLFW_KEY_MINUS
		KeyPeriod       Key = C.GLFW_KEY_PERIOD
		KeySlash        Key = C.GLFW_KEY_SLASH
		Key0            Key = C.GLFW_KEY_0
		Key1            Key = C.GLFW_KEY_1
		Key2            Key = C.GLFW_KEY_2
		Key3            Key = C.GLFW_KEY_3
		Key4            Key = C.GLFW_KEY_4
		Key5            Key = C.GLFW_KEY_5
		Key6            Key = C.GLFW_KEY_6
		Key7            Key = C.GLFW_KEY_7
		Key8            Key = C.GLFW_KEY_8
		Key9            Key = C.GLFW_KEY_9
		KeySemicolon    Key = C.GLFW_KEY_SEMICOLON
		KeyEqual        Key = C.GLFW_KEY_EQUAL
		KeyA            Key = C.GLFW_KEY_A
		KeyB            Key = C.GLFW_KEY_B
		KeyC            Key = C.GLFW_KEY_C
		KeyD            Key = C.GLFW_KEY_D
		KeyE            Key = C.GLFW_KEY_E
		KeyF            Key = C.GLFW_KEY_F
		KeyG            Key = C.GLFW_KEY_G
		KeyH            Key = C.GLFW_KEY_H
		KeyI            Key = C.GLFW_KEY_I
		KeyJ            Key = C.GLFW_KEY_J
		KeyK            Key = C.GLFW_KEY_K
		KeyL            Key = C.GLFW_KEY_L
		KeyM            Key = C.GLFW_KEY_M
		KeyN            Key = C.GLFW_KEY_N
		KeyO            Key = C.GLFW_KEY_O
		KeyP            Key = C.GLFW_KEY_P
		KeyQ            Key = C.GLFW_KEY_Q
		KeyR            Key = C.GLFW_KEY_R
		KeyS            Key = C.GLFW_KEY_S
		KeyT            Key = C.GLFW_KEY_T
		KeyU            Key = C.GLFW_KEY_U
		KeyV            Key = C.GLFW_KEY_V
		KeyW            Key = C.GLFW_KEY_W
		KeyX            Key = C.GLFW_KEY_X
		KeyY            Key = C.GLFW_KEY_Y
		KeyZ            Key = C.GLFW_KEY_Z
		KeyLeftBracket  Key = C.GLFW_KEY_LEFT_BRACKET
		KeyBackslash    Key = C.GLFW_KEY_BACKSLASH
		KeyRightBracket Key = C.GLFW_KEY_RIGHT_BRACKET
		KeyGraveAccent  Key = C.GLFW_KEY_GRAVE_ACCENT
		KeyWorld1       Key = C.GLFW_KEY_WORLD_1
		KeyWorld2       Key = C.GLFW_KEY_WORLD_2
		KeyEscape       Key = C.GLFW_KEY_ESCAPE
		KeyEnter        Key = C.GLFW_KEY_ENTER
		KeyTab          Key = C.GLFW_KEY_TAB
		KeyBackspace    Key = C.GLFW_KEY_BACKSPACE
		KeyInsert       Key = C.GLFW_KEY_INSERT
		KeyDelete       Key = C.GLFW_KEY_DELETE
		KeyRight        Key = C.GLFW_KEY_RIGHT
		KeyLeft         Key = C.GLFW_KEY_LEFT
		KeyDown         Key = C.GLFW_KEY_DOWN
		KeyUp           Key = C.GLFW_KEY_UP
		KeyPageUp       Key = C.GLFW_KEY_PAGE_UP
		KeyPageDown     Key = C.GLFW_KEY_PAGE_DOWN
		KeyHome         Key = C.GLFW_KEY_HOME
		KeyEnd          Key = C.GLFW_KEY_END
		KeyCapsLock     Key = C.GLFW_KEY_CAPS_LOCK
		KeyScrollLock   Key = C.GLFW_KEY_SCROLL_LOCK
		KeyNumLock      Key = C.GLFW_KEY_NUM_LOCK
		KeyPrintScreen  Key = C.GLFW_KEY_PRINT_SCREEN
		KeyPause        Key = C.GLFW_KEY_PAUSE
		KeyF1           Key = C.GLFW_KEY_F1
		KeyF2           Key = C.GLFW_KEY_F2
		KeyF3           Key = C.GLFW_KEY_F3
		KeyF4           Key = C.GLFW_KEY_F4
		KeyF5           Key = C.GLFW_KEY_F5
		KeyF6           Key = C.GLFW_KEY_F6
		KeyF7           Key = C.GLFW_KEY_F7
		KeyF8           Key = C.GLFW_KEY_F8
		KeyF9           Key = C.GLFW_KEY_F9
		KeyF10          Key = C.GLFW_KEY_F10
		KeyF11          Key = C.GLFW_KEY_F11
		KeyF12          Key = C.GLFW_KEY_F12
		KeyF13          Key = C.GLFW_KEY_F13
		KeyF14          Key = C.GLFW_KEY_F14
		KeyF15          Key = C.GLFW_KEY_F15
		KeyF16          Key = C.GLFW_KEY_F16
		KeyF17          Key = C.GLFW_KEY_F17
		KeyF18          Key = C.GLFW_KEY_F18
		KeyF19          Key = C.GLFW_KEY_F19
		KeyF20          Key = C.GLFW_KEY_F20
		KeyF21          Key = C.GLFW_KEY_F21
		KeyF22          Key = C.GLFW_KEY_F22
		KeyF23          Key = C.GLFW_KEY_F23
		KeyF24          Key = C.GLFW_KEY_F24
		KeyF25          Key = C.GLFW_KEY_F25
		KeyKP0          Key = C.GLFW_KEY_KP_0
		KeyKP1          Key = C.GLFW_KEY_KP_1
		KeyKP2          Key = C.GLFW_KEY_KP_2
		KeyKP3          Key = C.GLFW_KEY_KP_3
		KeyKP4          Key = C.GLFW_KEY_KP_4
		KeyKP5          Key = C.GLFW_KEY_KP_5
		KeyKP6          Key = C.GLFW_KEY_KP_6
		KeyKP7          Key = C.GLFW_KEY_KP_7
		KeyKP8          Key = C.GLFW_KEY_KP_8
		KeyKP9          Key = C.GLFW_KEY_KP_9
		KeyKPDecimal    Key = C.GLFW_KEY_KP_DECIMAL
		KeyKPDivide     Key = C.GLFW_KEY_KP_DIVIDE
		KeyKPMultiply   Key = C.GLFW_KEY_KP_MULTIPLY
		KeyKPSubtract   Key = C.GLFW_KEY_KP_SUBTRACT
		KeyKPAdd        Key = C.GLFW_KEY_KP_ADD
		KeyKPEnter      Key = C.GLFW_KEY_KP_ENTER
		KeyKPEqual      Key = C.GLFW_KEY_KP_EQUAL
		KeyLeftShift    Key = C.GLFW_KEY_LEFT_SHIFT
		KeyLeftControl  Key = C.GLFW_KEY_LEFT_CONTROL
		KeyLeftAlt      Key = C.GLFW_KEY_LEFT_ALT
		KeyLeftSuper    Key = C.GLFW_KEY_LEFT_SUPER
		KeyRightShift   Key = C.GLFW_KEY_RIGHT_SHIFT
		KeyRightControl Key = C.GLFW_KEY_RIGHT_CONTROL
		KeyRightAlt     Key = C.GLFW_KEY_RIGHT_ALT
		KeyRightSuper   Key = C.GLFW_KEY_RIGHT_SUPER
		KeyMenu         Key = C.GLFW_KEY_MENU
		KeyLast         Key = C.GLFW_KEY_LAST
	*/
	//Oh god
	keyInverse = make(map[int]Key)
	for _, v := range keyStore {
		keyInverse[v.KeyScancode] = v
	}
	keyListenerLock.Lock()
	keyListeners = make(map[string](map[int]chan int))
	for _, v := range keyStore {
		var chans map[int]chan int
		keyListeners[v.KeyName] = chans
	}
	keyListenerLock.Unlock()
}
