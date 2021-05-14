package update

import (
	"sync"
	"time"

	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/linmath"
	"github.com/jevans40/psychic-spork/objects"
)

var _ objects.ObjEnviroment = (*UpdateWorker)(nil)

type UpdateWorker struct {
	update      chan int
	render      chan struct{}
	event       chan []event.UpdateEvent
	coordinator chan []event.UpdateEvent
	waitGroup   *sync.WaitGroup
	toRender    []float32

	Objects map[int]objects.Object
	//This needs to be changed out for a new datastructure that keeps track of what
	eventQueue []event.UpdateEvent
	toSend     []event.UpdateEvent
}

func UpdateWorkerFactory(update chan int, render chan struct{}, eventm chan []event.UpdateEvent, corrdinatorEvent chan []event.UpdateEvent, wait *sync.WaitGroup) *UpdateWorker {
	ob := make(map[int]objects.Object)
	var eventq []event.UpdateEvent
	var send []event.UpdateEvent
	var toRender []float32
	return &UpdateWorker{update, render, eventm, corrdinatorEvent, wait, toRender, ob, eventq, send}
}

func (w *UpdateWorker) Start(renderchan chan []float32) {
	total := time.Duration(0)
	times := 0
	for {
		select {

		case <-w.render:
			start := time.Now()
			w.Render(renderchan)
			total = total + time.Since(start)
			times = times + 1
			if times%120 == 0 {
				//fmt.Printf("%v average per frame\n", total/120)
				total = time.Duration(0)
			}
		case updates := <-w.update:
			w.ProcessEvents()
			w.UpdateObjects(updates)
			w.PassEvents()
		case e := <-w.event:
			w.StoreEvents(e)
		}
	}
}

func (w *UpdateWorker) Render(renderchan chan []float32) {
	defer w.waitGroup.Done()
	i := 0
	for _, v := range w.Objects {
		v.Render(w.toRender[i*28 : (i+1)*28])
		i++
	}
	renderchan <- w.toRender[:]
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
			w.Objects[e.Receiver].SetEnviroment(w)
			w.Objects[e.Receiver].SendEvent(e)
			w.toRender = append(w.toRender, make([]float32, 28)...)
		} else if e.EventCode == event.UpdateEvent_RemoveObject {
			delete(w.Objects, e.Receiver)
			w.toRender = w.toRender[0 : len(w.toRender)-28]
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

func (w *UpdateWorker) GetWindowSize() linmath.PSPoint {
	return linmath.NewPSPoint(1080, 920)
}

func (w *UpdateWorker) GetEventCallback() objects.EventCallback {
	return w.SendEvent
}

func (w *UpdateWorker) GetMousePosition() linmath.PSPoint {
	return linmath.NewPSPoint(1080/2, 920/2)
}
