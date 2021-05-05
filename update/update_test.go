//+build v
package update

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/jevans40/psychic-spork/event"
	"github.com/jevans40/psychic-spork/objects"
)

//This is a really simple object that calculates the nth fibbo number recursively once per update.
//Once it reaches f(n) it creates a copy of itself with n set to some random number between 1 and 50 and starts again from 1.
//Once every n calculations it sends a message to its child node
type FibboObject struct {
	callback objects.EventCallback
	//Goal
	FibboN int
	//Current Number
	Iterations int
	//RandGenerator
	rand *rand.Rand

	//Generation
	gen      int
	messager int
	tosend   []int
}

func (f *FibboObject) Update(iteration int) {
	for v := 0; v < 10; v++ {
		Fibbonum := 1

		Last := 1
		for i := 0; i < f.FibboN; i++ {
			tmp := Fibbonum
			Fibbonum += Last
			Last = tmp
		}
	}
	f.callback(event.UpdateEvent{EventCode: event.UpdateEvent_PassMessage,
		Sender:   f.messager,
		Receiver: f.messager,
		Event:    event.UpdateEvent_PassMessageEvent{f.tosend}})
	f.Iterations++
	if f.Iterations == 10 && f.gen < 18 {
		//fmt.Printf("Generation %v has calculated %v Fibbonacci numbers, now creating generation %v\n", f.gen, f.FibboN, f.gen+1)
		r := rand.New(rand.NewSource(int64(iteration)))
		newFibbo := FibboObject{nil, 30, 0, r, f.gen + 1, 0, make([]int, 640)}
		f.gen = f.gen + 1
		f.callback(event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
			Sender:   f.messager,
			Receiver: -1,
			Event:    event.UpdateEvent_NewObjectEvent{&newFibbo}})
		f.FibboN = 0
		f.Iterations = 0
	}
}

func (f *FibboObject) SendEvent(e event.UpdateEvent) {
	//Eat Event for Now
	if e.EventCode == event.UpdateEvent_NewObject {
		if e.Receiver > 0 {
			f.messager = e.Receiver
		}
	}
	Fibbonum := 1
	Last := 1
	for i := 0; i < f.FibboN; i++ {
		tmp := Fibbonum
		Fibbonum += Last
		Last = tmp
	}
}

func (f *FibboObject) SetEventCallback(call objects.EventCallback) {
	f.callback = call
}

func (f *FibboObject) Render() []float32 {
	return []float32{}
}

func Clock(iterations int, clockChannel chan int, waitgroup *sync.WaitGroup) {
	lasttime := time.Now()
	for i := 0; i < iterations; i++ {
		//time.Sleep(time.Millisecond * 5)
		//UpdateCount := fmt.Sprintf("Update number %v", i)
		//myFigure := figure.NewFigure(UpdateCount, "", true)
		//myFigure.Print()
		clockChannel <- 1
		if i%10 == 1 {
			fmt.Printf("Update #%v took %v ms\n", i, time.Since(lasttime))
		}
		lasttime = time.Now()
	}
	waitgroup.Done()
}

func TestUpdate(t *testing.T) {
	var _ objects.Object = (*FibboObject)(nil)
	seed := 1337
	test1EventChan := make(chan []event.UpdateEvent)
	test1RenderChan := make(chan []float32)
	r := rand.New(rand.NewSource(int64(seed)))
	newFibbo := FibboObject{nil, r.Intn(31), 0, r, 0, 0, make([]int, 1024)}
	wait := sync.WaitGroup{}
	wait.Add(1)

	coordinator := CoordinatorFactory(test1EventChan, test1RenderChan)

	clockChannel := make(chan int)
	go Clock(5000, clockChannel, &wait)
	go coordinator.Start(clockChannel)
	test1EventChan <- []event.UpdateEvent{event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
		Receiver: -1,
		Sender:   -1,
		Event:    event.UpdateEvent_NewObjectEvent{&newFibbo}}}
	wait.Wait()
}

func BenchmarkUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var _ objects.Object = (*FibboObject)(nil)
		seed := 1337
		test1EventChan := make(chan []event.UpdateEvent)
		test1RenderChan := make(chan []float32)
		r := rand.New(rand.NewSource(int64(seed)))
		newFibbo := FibboObject{nil, 30, 0, r, 0, 0, make([]int, 1024)}
		wait := sync.WaitGroup{}
		wait.Add(1)

		coordinator := CoordinatorFactory(test1EventChan, test1RenderChan)

		clockChannel := make(chan int)
		go Clock(5000, clockChannel, &wait)
		go coordinator.Start(clockChannel)
		test1EventChan <- []event.UpdateEvent{event.UpdateEvent{EventCode: event.UpdateEvent_NewObject,
			Sender:   -1,
			Receiver: -1,
			Event:    event.UpdateEvent_NewObjectEvent{&newFibbo}}}
		wait.Wait()
	}
}

func BenchmarkFibbo(b *testing.B) {

}
