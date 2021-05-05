package update

import (
	"fmt"
	"runtime"
	"sync"

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
//If transferto is not -1, then the number represents the worker this object should be moved to
func (t *objectTracker) GetWorker(from int) (destination int, transferto int) {
	t.bounceCount++
	transferto = -1
	//After 100 messages, see if we should move to a closer worker (May need to be adjusted)
	if t.bounceCount > 100 {
		max := 0
		maxworker := -1
		avg := 0.0
		for i, v := range t.bounceArray {
			avg += float64(v)
			if v > max {
				max = v
				maxworker = i
			}
		}
		avg = avg / float64(len(t.bounceArray))

		//If there is alot of events coming from one worker then we should move this object to that worker
		if float64(max) > avg && maxworker != t.Workernum {
			t.Workernum = maxworker
			transferto = maxworker
			t.bounceCount = 0
		}
	}
	return t.Workernum, transferto

}

func CoordinatorFactory(eventChannel chan []event.UpdateEvent, renderChannel chan []float32) *Coordinator {
	var work []*UpdateWorker
	var chans []workerCommunicator
	object := make(map[int]objectTracker)
	var event []event.UpdateEvent
	wg := sync.WaitGroup{}
	cord := Coordinator{work, chans, object, event, 0, eventChannel, 0, &wg}
	cord.Init()
	return &cord

}

func (c *Coordinator) Init() {
	//Create Worker's communication channels
	numthreads := runtime.NumCPU() - 1
	for i := 0; i < numthreads; i++ {
		//Once tested add a buffer for better performance
		update := make(chan int)
		render := make(chan struct{})
		event := make(chan []event.UpdateEvent)
		c.workerChannels = append(c.workerChannels, workerCommunicator{UpdateChan: update,
			RenderChan: render,
			EventChan:  event})
	}

	//Create Workers
	for i := 0; i < len(c.workerChannels); i++ {
		c.workers = append(c.workers, UpdateWorkerFactory(c.workerChannels[i].UpdateChan, c.workerChannels[i].RenderChan, c.workerChannels[i].EventChan, c.eventsChannel, c.workerWait))
	}
}

func (c *Coordinator) Start(ClockChannel chan int) {
	for _, worker := range c.workers {
		go worker.Start()
	}
	for {
		Command := <-ClockChannel

		if Command == 1 {
			//Update
			//#1, Process stored events
			//Later this should go in the order of object numbers that sent the event.
			//For now random order is fine.
			//Also I can move alot of these to the workers so coordinator is just sync and timing

			//Create an array to store messages for all workers
			toMessage := make([][]event.UpdateEvent, runtime.NumCPU()-1)
			for _, v := range c.EventList {
				if v.EventCode == event.UpdateEvent_NewObject {
					ev := (v.Event).(event.UpdateEvent_NewObjectEvent)
					//Pick a random worker to give the object to
					selectedWorker := c.numObjects % len(c.workers)
					//Add it to the object list
					c.objects[c.numObjects] = objectTracker{selectedWorker, make([]int, len(c.workers)), 0}
					//Set Object number in the events
					v.Receiver = c.numObjects
					if c.numObjects%100 == 0 {
						fmt.Printf("Adding object #%v\n", c.numObjects)
					}
					//Pass Event to the worker
					//c.workerChannels[selectedWorker].EventChan <- Event{event.UpdateEvent_NewObject, ev}
					toMessage[selectedWorker] = append(toMessage[selectedWorker],
						event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
							Sender:   -1,
							Receiver: v.Receiver,
							Event:    ev})
					c.numObjects++
				} else if v.EventCode == event.UpdateEvent_RemoveObject {
					//Remove it from the object list
					_, ok := c.objects[v.Receiver]
					if ok {
						delete(c.objects, v.Receiver)
					}
					//I dont think anything else needs to be done actually
				} else if v.EventCode == event.UpdateEvent_PassMessage {
					ev := (v.Event).(event.UpdateEvent_PassMessageEvent)
					//Verify that the receiver exists
					_, ok := c.objects[v.Receiver]
					if !ok {
						//Notify sender that message was not received
						//c.workerChannels[c.objects[ev.Sender].Workernum].EventChan <- Event{event.UpdateEvent_FailedSendMessage, event.UpdateEvent_FailedSendMessageEvent{ev.Sender, ev.Message}}
						toMessage[c.objects[v.Sender].Workernum] = append(toMessage[c.objects[v.Sender].Workernum],
							event.UpdateEvent{EventCode: event.UpdateEvent_FailedSendMessage,
								Sender:   v.Sender,
								Receiver: v.Sender,
								Event:    event.UpdateEvent_FailedSendMessageEvent{ev.Message}})
						continue
					}
					//Else just pass the message along
					//c.workerChannels[c.objects[ev.ObjectNumber].Workernum].EventChan <- Event{event.UpdateEvent_PassMessage, ev}
					toMessage[c.objects[v.Sender].Workernum] = append(toMessage[c.objects[v.Sender].Workernum],
						event.UpdateEvent{EventCode: event.UpdateEvent_PassMessage,
							Sender:   v.Sender,
							Receiver: v.Receiver,
							Event:    event.UpdateEvent_PassMessageEvent{ev.Message}})
				}
			}

			//Send collected messages
			for i, v := range toMessage {
				if v != nil {
					c.workerChannels[i].EventChan <- v
				}
			}

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
			continue
		}

		if Command == 2 {
			c.workerWait.Add(len(c.workerChannels))
			for i, v := range c.workerChannels {
				_ = i
				v.RenderChan <- struct{}{}
			}

		}
		if Command == 3 {
			break
		}
	}
}

/*
func (cd *coordinator) StartCoordinator(renderer chan event.AsyncEvent, numWorkers int) {
	var wg sync.WaitGroup
	var replies chan event.AsyncEvent
	for i := 0; i < numWorkers; i++ {
		workerchan := make(chan event.AsyncEvent) //For performance reasons please make this a buffered channel later when I'm sure this wont die
		worker := WorkerFactory(&wg, replies)
		worker.AddListener(workerchan)
		cd.workers = append(cd.workers, worker)
		cd.workerSendChannel = append(cd.workerSendChannel, workerchan)
	}
	cd.entities = make(map[int]int)
	for i := 0; i < numWorkers; i++ {
		go cd.workers[i].StartUpdateLoop()
	}
	wg.Add(1)
	for {
		cycleStartTime := time.Now()
		//Process requests
		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {

			cd.workerSendChannel[i] <- event.UpdateEventFactory(cycleStartTime)
		}
		wg.Wait()
		for {
			select {
			case i := <-replies:
				switch v := i.(type) {
				case *event.EReplyEvent:
					cd.workerSendChannel[cd.entities[v.Entity]] <- i
				default:
					logging.Log.Critical("ERROR UNKNOWN EVENT RECIEVED IN THE COORDINATOR REPLY SECTION")
				}
			default:
				goto FINISHEDREPLIESLABEL
			}

		}
	FINISHEDREPLIESLABEL:
		wg.Add(numWorkers)
		for i := 0; i < numWorkers; i++ {
			cd.workerSendChannel[i] <- event.DeltaEventFactory()
		}
		wg.Wait()

		for {
			select {
			case i := <-replies:
				switch v := i.(type) {
				case *event.EDeltaEvent:
					for j := 0; j < 28; j++ {
						cd.renderList[v.Entitynum*28+j] = v.Deltas[j]
					}
				default:
					logging.Log.Critical("ERROR UNKNOWN EVENT RECIEVED IN THE COORDINATOR REPLY SECTION")
				}
			default:
				goto FINISHEDDELTASLABEL
			}

		}
	FINISHEDDELTASLABEL:
		//Sync workers
	}
}
*/
