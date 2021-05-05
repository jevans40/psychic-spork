package update

import (
	"sync"

	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/objects"
)

type UpdateWorker struct {
	update      chan int
	render      chan struct{}
	event       chan []event.UpdateEvent
	coordinator chan []event.UpdateEvent
	waitGroup   *sync.WaitGroup

	Objects map[int]objects.Object
	//This needs to be changed out for a new datastructure that keeps track of what
	eventQueue []event.UpdateEvent
	toSend     []event.UpdateEvent
}

func UpdateWorkerFactory(update chan int, render chan struct{}, eventm chan []event.UpdateEvent, corrdinatorEvent chan []event.UpdateEvent, wait *sync.WaitGroup) *UpdateWorker {
	ob := make(map[int]objects.Object)
	var eventq []event.UpdateEvent
	var send []event.UpdateEvent
	return &UpdateWorker{update, render, eventm, corrdinatorEvent, wait, ob, eventq, send}
}

func (w *UpdateWorker) Start() {
	for {
		select {

		case <-w.render:
			w.Render()
		case updates := <-w.update:
			w.ProcessEvents()
			w.UpdateObjects(updates)
			w.PassEvents()
		case e := <-w.event:
			w.StoreEvents(e)
		}
	}
}

func (w *UpdateWorker) Render() {
	defer w.waitGroup.Done()
	renderlist := make([]float32, 28*len(w.Objects))
	for i, v := range w.Objects {
		copy(renderlist[i*28:(i+1)*28], v.Render()[:])
	}
}

func (w *UpdateWorker) UpdateObjects(time int) {
	for _, v := range w.Objects {
		v.Update(time)
	}
}

func (w *UpdateWorker) StoreEvents(e []event.UpdateEvent) {
	w.eventQueue = append(w.eventQueue, e...)
}

func (w *UpdateWorker) ProcessEvents() {
	for _, e := range w.eventQueue {

		if e.EventCode == event.UpdateEvent_NewObject {
			ev := (e.Event).(event.UpdateEvent_NewObjectEvent)
			w.Objects[e.Receiver] = (ev.Object).(objects.Object)
			w.Objects[e.Receiver].SetEventCallback(w.SendEvent)
			w.Objects[e.Receiver].SendEvent(e)
		} else if e.EventCode == event.UpdateEvent_RemoveObject {
			delete(w.Objects, e.Receiver)
		} else if e.EventCode == event.UpdateEvent_PassMessage {
			w.Objects[e.Receiver].SendEvent(e)
		} else if e.EventCode == event.UpdateEvent_FailedSendMessage {
			w.Objects[e.Receiver].SendEvent(e)
		}
	}
	w.eventQueue = w.eventQueue[:0]
}

func (w *UpdateWorker) SendEvent(e event.UpdateEvent) {

	if e.EventCode == event.UpdateEvent_PassMessage {
		_, ok := w.Objects[e.Receiver]
		if ok {
			w.eventQueue = append(w.eventQueue, e)
			return
		}
	} else if e.EventCode == event.UpdateEvent_FailedSendMessage {
		_, ok := w.Objects[e.Receiver]
		if ok {
			w.eventQueue = append(w.eventQueue, e)
			return
		}
	}
	w.toSend = append(w.toSend, e)
}

func (w *UpdateWorker) PassEvents() {
	defer w.waitGroup.Done()
	if len(w.toSend) != 0 {

		w.coordinator <- w.toSend
	}
	w.toSend = w.toSend[:0]
}

/*
type UpdateWorker interface {
	StartUpdateLoop()
	AddListener(listener chan event.AsyncEvent)
}

type simpleWorker struct {
	messagingChan <-chan event.AsyncEvent
	replyChan     chan event.AsyncEvent
	waitG         *sync.WaitGroup
	entityList    Entity
}

type EventList interface {
}

func (sw *simpleWorker) StartUpdateLoop() {
	closed := false
	for !closed {
		select {
		case i, ok := <-sw.messagingChan:
			//Check if chanel was closed
			if !ok {
				closed = true
			}

			//Update entities
			switch v := i.GetType(); v {
			case event.UpdateEvent:

			}
		}
	}
}

func (sw *simpleWorker) AddListener(listener chan event.AsyncEvent) {
	sw.messagingChan = listener
}

func WorkerFactory(wg *sync.WaitGroup, replychan chan event.AsyncEvent) UpdateWorker {
	worker := simpleWorker{waitG: wg, replyChan: replychan}
	return &worker
}
*/
