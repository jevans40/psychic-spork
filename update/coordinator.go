package update

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/jevans40/psychic-spork/event"
)

type Coordinator struct {
	workers        []*UpdateWorker
	workerChannels []workerCommunicator
	objects        map[int]objectTracker
	EventList      []event.UpdateEvent
	numObjects     int
	eventsChannel  chan []event.UpdateEvent
	updates        int
	workerWait     *sync.WaitGroup
	renderer       chan []float32
	toMessage      [][]event.UpdateEvent

	gameEventsChannel chan event.UpdateEvent
}

type workerCommunicator struct {
	//Communication on this channel indicate an update signal
	//int: The number of updates so far time is determined by clock thread
	UpdateChan chan int

	//Communication on this channel indicate that the worker should send updated verticies when it's done with its current task.
	//struct is a zero byte structure that just indicates a message
	RenderChan chan struct{}

	//
	EventChan chan []event.UpdateEvent
}

type objectTracker struct {
	//-1 here indicates that the object has been delted
	Workernum   int
	bounceArray []int
	bounceCount int
}

//Get the worker this object is currently on
//REFACTOR
func (c *Coordinator) AddMissMessage(from int, to int) {
	fr := c.objects[from]
	fr.bounceCount++
	worker1 := c.objects[from].Workernum
	tr := c.objects[to]
	tr.bounceCount++
	worker2 := c.objects[from].Workernum

	//After 100 messages, see if we should move to a closer worker (May need to be adjusted)
	if fr.bounceCount > 100 {
		c.toMessage[worker1] = append(c.toMessage[worker1], *event.NewGetMessageLogsFactory().To(from).From(-1))
	}
	if tr.bounceCount > 100 {
		c.toMessage[worker2] = append(c.toMessage[worker2], *event.NewGetMessageLogsFactory().To(to).From(-1))
	}
}

//REFACTOR
func CoordinatorFactory(eventChannel chan []event.UpdateEvent, renderChannel chan []float32) *Coordinator {
	var work []*UpdateWorker
	var chans []workerCommunicator
	object := make(map[int]objectTracker)
	var event []event.UpdateEvent
	wg := sync.WaitGroup{}
	cord := Coordinator{work, chans, object, event, 0, eventChannel, 0, &wg, renderChannel, nil, nil}
	cord.Init()
	return &cord

}

//REFACTOR
func (c *Coordinator) Init() {
	//Create Worker's communication channels
	numthreads := 6 //runtime.NumCPU() - 1
	for i := 0; i < numthreads; i++ {
		//Once tested add a buffer for better performance
		update := make(chan int)
		render := make(chan struct{})
		event := make(chan []event.UpdateEvent)
		c.workerChannels = append(c.workerChannels, workerCommunicator{UpdateChan: update,
			RenderChan: render,
			EventChan:  event})
	}
	c.toMessage = make([][]event.UpdateEvent, 0, len(c.workerChannels))
	//Create Workers
	for i := 0; i < len(c.workerChannels); i++ {
		c.workers = append(c.workers, UpdateWorkerFactory(c.workerChannels[i].UpdateChan, c.workerChannels[i].RenderChan, c.workerChannels[i].EventChan, c.eventsChannel, c.workerWait))
		c.toMessage = append(c.toMessage, []event.UpdateEvent{})
	}
}

//REFACTOR
func (c *Coordinator) Start(ClockChannel chan int, gameEventsChannel chan event.UpdateEvent) {
	c.gameEventsChannel = gameEventsChannel
	fromWorkerRender := make(chan []float32, runtime.NumCPU())

	for _, worker := range c.workers {
		go worker.Start(fromWorkerRender)
	}
	updat := 0
	lastPrintUpdate := time.Duration(0)
	rend := 0
	lastPrintRender := time.Duration(0)

	for {
		Command := <-ClockChannel

		if Command == 1 {
			//Update
			//#1, Process stored events
			//Later this should go in the order of object numbers that sent the event.
			//For now random order is fine.
			//Also I can move alot of these to the workers so coordinator is just sync and timing

			//Create an array to store messages for all workers
			start := time.Now()
			for i, _ := range c.toMessage {
				c.toMessage[i] = c.toMessage[i][:0]
			}
			for _, v := range c.EventList {
				switch v.EventCode {

				case event.UpdateEvent_NewObject:
					c.HandleNewObjectEvent(v)
				case event.UpdateEvent_RemoveObject:
					c.HandleRemoveObjectEvent(v)
				case event.UpdateEvent_PassMessage:
					c.HandlePassMessageEvent(v)
				case event.SubscribeEvent_WindowResize:
					c.HandleSubscribeEvent(v)
				}
			}

			//Send collected messages
			for i, v := range c.toMessage[0 : len(c.toMessage)-1] {
				if v != nil {
					c.workerChannels[i].EventChan <- v
				}
			}

			//eventLoopChannel <- toMessage[len(toMessage)-1]

			c.EventList = make([]event.UpdateEvent, 0)
			//#2 Send Update signal
			c.updates++
			c.workerWait.Add(len(c.workerChannels))
			for i, v := range c.workerChannels {
				_ = i
				v.UpdateChan <- c.updates
			}
			//#3 Collect Events and wait for all workers to finish
			done := make(chan struct{})

			go func() {
				c.workerWait.Wait()
				done <- struct{}{}
			}()

			stop := false
			for {
				select {
				case x := <-c.eventsChannel:
					c.EventList = append(c.EventList, x...)
				case <-done:
					//Keep reading in events, but now if we dont get any we know that we can stop
					stop = true
				default:
					if stop {
						goto END_EVENT_GET
					}
				}
			}
		END_EVENT_GET:
			updat++
			lastPrintUpdate = lastPrintUpdate + time.Since(start)
			if updat%120 == 0 {
				fmt.Printf("%v average Update time\n", lastPrintUpdate/120)
				lastPrintUpdate = time.Duration(0)
			}
			continue
		}

		if Command == 2 {
			c.workerWait.Add(len(c.workerChannels))
			for i, v := range c.workerChannels {
				_ = i
				v.RenderChan <- struct{}{}
			}
			var toRender []float32 = make([]float32, c.numObjects*28)
			done := make(chan struct{})
			go func() {
				c.workerWait.Wait()
				done <- struct{}{}
			}()
			total := 0
			stop := false
			for {
				select {
				case x := <-fromWorkerRender:
					start := time.Now()
					if x != nil {
						copy(toRender[total:total+len(x)-1], x[:])
					}
					total = total + len(x)
					rend++
					lastPrintRender = lastPrintRender + time.Since(start)
					if rend%(120*len(c.workerChannels)) == 0 {
						fmt.Printf("%v average Render time\n", lastPrintRender/120)
						lastPrintRender = time.Duration(0)
					}
				case <-done:
					//Keep reading in events, but now if we dont get any we know that we can stop
					stop = true
				default:
					if stop {
						goto END_RENDER_GET
					}
				}
			}
		END_RENDER_GET:
			c.renderer <- toRender
			continue
		}
		if Command == 3 {
			break
		}
	}
}

//----------------------------------------------------------
//Here are handlers for events the coordinator will receive
//----------------------------------------------------------

//Coordinator should create a new objectTracker and pass the event on to a random worker
func (c *Coordinator) HandleNewObjectEvent(e event.UpdateEvent) {
	ev := (e.Event).(event.UpdateEvent_NewObjectEvent)
	//Pick a random worker to give the object to
	selectedWorker := c.numObjects % len(c.workers)
	//Add it to the object list
	c.objects[c.numObjects] = objectTracker{selectedWorker, make([]int, len(c.workers)), 0}
	//Set Object number in the events
	e.Receiver = c.numObjects
	//Pass Event to the worker
	//c.workerChannels[selectedWorker].EventChan <- Event{event.UpdateEvent_NewObject, ev}
	c.toMessage[selectedWorker] = append(c.toMessage[selectedWorker],
		event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
			Sender:   -1,
			Receiver: e.Receiver,
			Event:    ev})
	c.numObjects++
}

//Coordinator should remove objecttracker, but is not responsible for the object being removed
//Does nothing if the sender is not the receiver
func (c *Coordinator) HandleRemoveObjectEvent(e event.UpdateEvent) {
	if e.Sender != e.Receiver {
		return
	}
	_, ok := c.objects[e.Receiver]
	if ok {
		delete(c.objects, e.Receiver)
	}
}

func (c *Coordinator) HandlePassMessageEvent(e event.UpdateEvent) {
	ev := (e.Event).(event.UpdateEvent_PassMessageEvent)
	//Verify that the receiver exists
	_, ok := c.objects[e.Receiver]
	if !ok {
		//Notify sender that message was not received
		//c.workerChannels[c.objects[ev.Sender].Workernum].EventChan <- Event{event.UpdateEvent_FailedSendMessage, event.UpdateEvent_FailedSendMessageEvent{ev.Sender, ev.Message}}
		c.toMessage[c.objects[e.Sender].Workernum] = append(c.toMessage[c.objects[e.Sender].Workernum],
			*event.NewFailedSendMessage(ev.Message).To(e.Sender).From(e.Sender))
		return
	}
	//Else just pass the message along
	//c.workerChannels[c.objects[ev.ObjectNumber].Workernum].EventChan <- Event{event.UpdateEvent_PassMessage, ev}
	c.toMessage[c.objects[e.Receiver].Workernum] = append(c.toMessage[c.objects[e.Receiver].Workernum], e)
}

func (c *Coordinator) HandleFailedSendMessageEvent(e event.UpdateEvent) {
	_, ok := c.objects[e.Receiver]
	if !ok {
		return
	}
	c.toMessage[c.objects[e.Receiver].Workernum] = append(c.toMessage[c.objects[e.Receiver].Workernum], e)
}
func (c *Coordinator) HandleSubscribeEvent(e event.UpdateEvent) {
	c.gameEventsChannel <- e
}
func (c *Coordinator) HandleUnSubscribeEvent(e event.UpdateEvent) {
	c.gameEventsChannel <- e
}
func (c *Coordinator) HandleWindowResizeEvent(e event.UpdateEvent) {
	c.gameEventsChannel <- e
}
func (c *Coordinator) UnSubscribe_WindowResizeEvent(e event.UpdateEvent) {
	c.gameEventsChannel <- e
}
func (c *Coordinator) HandleGetMessageLogsEvent(event event.UpdateEvent) {

}

func (c *Coordinator) HandleReturnMessageLogs(event event.UpdateEvent) {

}
