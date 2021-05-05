package game

import "time"

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
	for {
		if time.Since(lastRender) > c.renderTickDur {
			c.commChan <- 2
		}
		if time.Since(lastUpdate) > c.updateTickDur {
			c.commChan <- 1
		}
	}
}
