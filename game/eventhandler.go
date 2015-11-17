package game

import (
	"bufio"
	"fmt"
	"os"
)

// EventHandler object
// FIX:This struct is now useless, because the handler is not ignite.
type EventHandler struct {
	eventChs map[string]chan bool
	isHandle bool
}

// NewEventHandler is constructor of EventHandler
// @Param channels イベントハンドラ用チャネル
// return EventHandler
func NewEventHandler(channels map[string]chan bool) *EventHandler {
	return &EventHandler{
		eventChs: channels,
	}
}

// StartHandle starts reading input
func (eh *EventHandler) StartHandle() {
	scanner := bufio.NewScanner(os.Stdin)
	eh.isHandle = true
	for scanner.Scan() {
		if eh.isHandle {
			// 検知終了
			fmt.Println("Finished handler")
			break
		}
		fmt.Println("aaa")
		for k, v := range eh.eventChs {
			if scanner.Text() == k {
				v <- true
			}
		}
	}
}

// StopHandle stops the loop of handler
func (eh *EventHandler) StopHandle() {
	eh.isHandle = false
}

// CloseChs closes all channels
func (eh *EventHandler) CloseChs() {
	for _, v := range eh.eventChs {
		close(v)
	}
}
