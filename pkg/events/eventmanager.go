package events

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	eventQueuesCapacity                                       = 100000
	idleDispatcherSleepTime                     time.Duration = 5 * time.Millisecond
	registeringListenerWhileRunningErrorMessage               = "Tried to register listener while running event loop. Registering listeners is not thread safe therefore prohibited after starting event loop."
)

type EventIdType string

// Defines a game event
type Event interface {
	EventId() EventIdType
}

// Defines a dispatcher. Any service wishing to dispatch events should register
type EventDispatcher interface {
	RegisterEventTypes(*EventManager)
}

// Defines an event handler
type EventHandler interface {
	Handle(event Event)
	RunsInOwnThread() bool
}

type EventManager struct {
	running        bool
	quit           chan struct{}
	eventQueue     chan Event
	eventListeners map[EventIdType][]EventHandler
}

func NewEventManager() *EventManager {
	return &EventManager{
		running:        false,
		quit:           make(chan struct{}),
		eventQueue:     make(chan Event, 10000),
		eventListeners: make(map[EventIdType][]EventHandler),
	}
}

func (m *EventManager) Dispatch(event Event) error {
	select {
	case m.eventQueue <- event:
		return nil
	default:
		return fmt.Errorf("failed to dispatch event id %s", event.EventId())
	}
}

// Register an event type
func (m *EventManager) RegisterEventType(_id EventIdType) error {
	m.panicWhenEventLoopRunning()
	_, ok := m.eventListeners[_id]
	// If the key exists
	if ok {
		return fmt.Errorf("event %s already registered", _id)
	} else {
		m.eventListeners[_id] = []EventHandler{}
	}
	return nil
}

// Register a handler for an event type
func (m *EventManager) RegisterHandler(_id EventIdType, _h EventHandler) error {
	m.panicWhenEventLoopRunning()
	if _, ok := m.eventListeners[_id]; !ok {
		return fmt.Errorf("event %s has not been registered", _id)
	}
	m.eventListeners[_id] = append(m.eventListeners[_id], _h)
	return nil
}

// Run the event loop
func (m *EventManager) Run() {
	defer func() {
		m.running = false
		log.Info("[EVENT] Event manager quitting")
	}()
	m.running = true

	log.Info("[EVENT] Event manager running")

	for {
		select {

		case e := <-m.eventQueue:
			if listeners, ok := m.eventListeners[e.EventId()]; ok {
				log.Infof("[EVENT] %s", e.EventId())
				for _, handler := range listeners {
					if handler.RunsInOwnThread() {
						go handler.Handle(e)
					} else {
						handler.Handle(e)
					}
				}
			} else {
				log.Errorf("[EVENT] event %s not registered ", e.EventId())
			}
		case _, ok := <-m.quit:
			if !ok {
				return
			}
		default:
			time.Sleep(idleDispatcherSleepTime)
		}
	}
}

// Stop event handler
func (m *EventManager) Stop() {
	close(m.quit)
}

func (m *EventManager) panicWhenEventLoopRunning() {
	if m.running {
		panic(registeringListenerWhileRunningErrorMessage)
	}
}
