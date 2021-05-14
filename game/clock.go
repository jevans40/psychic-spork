package game

import (
	"time"
)

//A simple struct that can be configured to tick update and render loops at set intervals
//TODO: When generics come out update this clock to be a generic clock
type Clock struct {
	updateTickDur time.Duration
	renderTickDur time.Duration

	commChan chan int
}

func (c *Clock) SetDurations(renderDur time.Duration, updateTickDur time.Duration) {

}

func (c *Clock) SetChannels(newchan chan int) {
	c.commChan = newchan
}

func (c *Clock) Start() {
	lastUpdate := time.Now()
	lastRender := time.Now()
	numUpdate := 0
	lastPrintUpdate := time.Duration(0)
	numRender := 0
	lastPrintRender := time.Duration(0)
	for {
		if time.Since(lastRender) > c.renderTickDur {
			start := time.Now()
			c.commChan <- 2
			lastPrintRender = lastPrintRender + time.Since(start)
			numRender++
			if numRender%120 == 0 {
				//fmt.Printf("%v average render time\n", lastPrintRender/120)
				lastPrintRender = time.Duration(0)
			}

		}
		if time.Since(lastUpdate) > c.updateTickDur {
			start := time.Now()
			c.commChan <- 1
			lastPrintUpdate = lastPrintUpdate + time.Since(start)
			numUpdate++
			if numUpdate%120 == 0 {
				//fmt.Printf("%v average update time\n", lastPrintUpdate/120)
				lastPrintUpdate = time.Duration(0)
			}
		}
	}
}
