package main

import (
	"fmt"
	"log"
	"time"
)


const (
	LiftStateStop = iota
	LiftStateMove
)


type Info struct {
	State    int
	Floor    int
	IsMoveUp bool
}


type Lift struct {
	doorsDelay time.Duration
	floorDelay time.Duration
	infoCh     chan *Info
	lockCh     chan struct{}
}


func NewLift(count, height, speed uint, delay time.Duration) (*Lift, error) {
	if speed == 0 {
		return nil, fmt.Errorf("speed can't be zero")
	}
	if height == 0 {
		return nil, fmt.Errorf("floor height must be > 0")
	}
	if count < 5 || count > 10 {
		return nil, fmt.Errorf("floor count must be in range: from 5 to 10")
	}
	if speed > height {
		return nil, fmt.Errorf("speed can't be more than height, it's super lift")
	}

	return &Lift{
		doorsDelay: delay,
		floorDelay: time.Second * (time.Duration(height / speed)),
		infoCh:     make(chan *Info),
		lockCh:     make(chan struct{}),
	}, nil
}


func (l *Lift) Run() (<-chan *Info, chan<- struct{}) {
	return l.infoCh, l.lockCh
}


func (l *Lift) Move(from, to int) {
	if from == to {
		l.OpenCloseDoors()
		return
	}

	isMoveUp := from < to
	if isMoveUp {
		for i := from; i < to; i++ {
			time.Sleep(l.floorDelay)
			log.Printf("Lift move the floor: #%d\n", i)
			l.sendInfo(LiftStateMove, i, isMoveUp)
		}
	} else {
		for i := from; i > to; i-- {
			time.Sleep(l.floorDelay)
			log.Printf("Lift move the floor: #%d\n", i)
			l.sendInfo(LiftStateMove, i, isMoveUp)
		}
	}

	l.sendInfo(LiftStateStop, to, isMoveUp)
}


func (l *Lift) OpenCloseDoors() {
	log.Printf("Lift open doors\n")
	time.Sleep(l.doorsDelay)
	log.Printf("Lift close doors\n")
}


func (l *Lift) sendInfo(state, floor int, isMoveUp bool) {
	l.infoCh <- &Info{
		State:    state,
		Floor:    floor,
		IsMoveUp: isMoveUp,
	}


	<-l.lockCh
}

func main()  {

	}


