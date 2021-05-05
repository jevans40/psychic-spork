package event

type UpdateEvent struct {
	EventCode int
	Sender    int
	Receiver  int
	Event     interface{}
}

const UpdateEvent_NewObject = 0
const UpdateEvent_RemoveObject = 1
const UpdateEvent_PassMessage = 2
const UpdateEvent_FailedSendMessage = 3
const UpdateEvent_Subscribe = 4
const UpdateEvent_UnSubscribe = 5

//Event is sent with an included object to add to a random worker
//object number is assigned by the coordinator and will be ignored -1 by convention
type UpdateEvent_NewObjectEvent struct {
	//Once objects are better defined this should become more selective
	Object interface{}
}

//Simply updates the coordinator that the object has been deleted.
type UpdateEvent_RemoveObjectEvent struct {
}

//Sends a message from one object to another specified here
//Sender is the object ID of the object that sent the message
type UpdateEvent_PassMessageEvent struct {
	Message interface{}
}

type UpdateEvent_FailedSendMessageEvent struct {
	Message interface{}
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
//Closes the channel
type UpdateEvent_UnSubscribeEvent struct {
	ListeningChannel      chan int
	EventSubscriptionName string
}
