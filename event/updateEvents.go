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
