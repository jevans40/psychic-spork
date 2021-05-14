package event

import "github.com/jevans40/psychic-spork/linmath"

type UpdateEvent struct {
	EventCode int
	Sender    int
	Receiver  int
	Event     interface{}
}

func (u *UpdateEvent) To(Receiver int) *UpdateEvent {
	u.Receiver = Receiver
	return u
}

func (u *UpdateEvent) From(Sender int) *UpdateEvent {
	u.Sender = Sender
	return u
}

const UpdateEvent_NewObject = 0
const UpdateEvent_RemoveObject = 1
const UpdateEvent_PassMessage = 2
const UpdateEvent_FailedSendMessage = 3
const UpdateEvent_Subscribe = 4
const UpdateEvent_UnSubscribe = 5
const SubscribeEvent_WindowResize = 6
const SubscribeEvent_UnSubscribe_WindowResize = 7
const UpdateEvent_GetMessageLogs = 8
const UpdateEvent_ReturnMessageLogs = 9

//REFACTOR all events should have an assocated factory function

//Event is sent with an included object to add to a random worker
//object number is assigned by the coordinator and will be ignored -1 by convention
type UpdateEvent_NewObjectEvent struct {
	//Once objects are better defined this should become more selective
	Object interface{}
}

func NewObjectEventFactory(object interface{}) *UpdateEvent {
	return &UpdateEvent{EventCode: UpdateEvent_NewObject,
		Sender:   -1,
		Receiver: -1,
		Event:    &UpdateEvent_NewObjectEvent{},
	}

}

//Simply updates the coordinator that the object has been deleted.
//This should ONLY be sent by the object being deleted
type UpdateEvent_RemoveObjectEvent struct {
}

//Sends a message from one object to another specified here
//Sender is the object ID of the object that sent the message
type UpdateEvent_PassMessageEvent struct {
	Message interface{}
}

func NewPassMessage(Message interface{}) *UpdateEvent {
	return &UpdateEvent{EventCode: UpdateEvent_PassMessage,
		Sender:   -1,
		Receiver: -1,
		Event:    UpdateEvent_PassMessageEvent{Message}}
}

type UpdateEvent_FailedSendMessageEvent struct {
	Message interface{}
}

func NewFailedSendMessage(Message interface{}) *UpdateEvent {
	return &UpdateEvent{EventCode: UpdateEvent_FailedSendMessage,
		Sender:   -1,
		Receiver: -1,
		Event:    UpdateEvent_FailedSendMessageEvent{Message}}
}

//This is sent to glfwEvents to subscribe to certain listeners
//Channels MUST be buffered, or they may never receive a message
//If the channel is full then the message will also be dropped
//Keeps track of the channel by sender name
type UpdateEvent_SubscribeEvent struct {
	ListeningChannel      chan int
	EventSubscriptionName string
}

//This is sent to glfwEvents to subscribe to certain listeners
//Unsubscribes from channel associated with the Sender
type UpdateEvent_UnSubscribeEvent struct {
	EventSubscriptionName string
}

//Adds a listener to the window resize event
type SubscribeEvent_WindowResizeEvent struct {
	ListeningChannel chan linmath.PSPoint
}

//Removes a listener from the window resize event
type SubscribeEvent_UnSubscribe_WindowResizeEvent struct {
}

type UpdateEvent_GetMessageLogsEvent struct {
}

func NewGetMessageLogsFactory() *UpdateEvent {
	return &UpdateEvent{EventCode: UpdateEvent_GetMessageLogs,
		Sender:   -1,
		Receiver: -1,
		Event:    &UpdateEvent_GetMessageLogsEvent{},
	}
}

type UpdateEvent_ReturnMessageLogsEvent struct {
	counts []int
}
