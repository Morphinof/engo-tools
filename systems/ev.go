package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"log"
	"tools/components"
)

type EventSystem struct {
	world  *ecs.World
	hears  map[engo.MessageHandlerId]*components.Event
	events map[string]*components.Event
}

func (es *EventSystem) New(w *ecs.World) {
	es.world = w
	es.hears = make(map[engo.MessageHandlerId]*components.Event, 0)
	es.events = make(map[string]*components.Event, 0)
}

func (es *EventSystem) Update(dt float32) {
	if engo.Input.Button("F2").JustPressed() {
		es.Debug()
	}
}

func (es *EventSystem) Remove(e ecs.BasicEntity) {
	for name, event := range es.events {
		if event.ID() == e.ID() {
			if event.MessageHandlerId > 0 {
				engo.Mailbox.StopListen(event.Name, event.MessageHandlerId)
				delete(es.hears, event.MessageHandlerId)
			}
			delete(es.events, name)
			fmt.Printf("ES:Remove #%d %p removed\n", e.ID(), &e)
			return
		}
	}

	log.Fatalf("ES:Remove Event %d %p not found!\n", e.ID(), &e)
}

func (es *EventSystem) NewEvent(name string) {
	if es.events[name] != nil {
		log.Fatalf("ES:NewEvent: Event %s already exists\n", name)
		return
	}

	e := components.Event{
		BasicEntity: ecs.NewBasic(),
		Name:        name,
		Active:      true,
		Data:        make(map[string]interface{}),
	}

	es.events[name] = &e
	fmt.Printf("ES:NewEvent: Created %d %s\n", e.ID(), e.Name)
}

func (es *EventSystem) Get(name string) *components.Event {
	if es.events[name] == nil {
		log.Fatalf("ES: Unknow event %s\n", name)
	}

	return es.events[name]
}

func (es *EventSystem) Listen(name string, handler engo.MessageHandler) {
	if es.events[name] == nil {
		log.Fatalf("ES:Listen: Unknow event %s\n", name)
	}

	e := es.events[name]
	handlerId := engo.Mailbox.Listen(e.Name, handler)
	fmt.Printf("ES:Listen: Listening event %p : %d %s - %d\n", e, e.ID(), e.Name, handlerId)
	es.hears[handlerId] = e
	e.Listened = true
}

func (es *EventSystem) Dispatch(name string, data map[string]any) {
	if es.events[name] == nil {
		log.Fatalf("ES:Dispatch: Unknow event %s\n", name)
	}

	e := es.events[name]
	e.Data = data
	engo.Mailbox.Dispatch(e)
	e.Dispatched++
	fmt.Printf("ES:Dispatch: %d %s dispatched\n", e.ID(), e.Name)
}

func (es *EventSystem) Disable(name string) {
	if es.events[name] == nil {
		log.Fatalf("ES:Disable: Unknow event %s\n", name)
	}

	e := es.events[name]
	hear := es.hears[e.MessageHandlerId]
	if hear == nil {
		return
	}
	hear.Active = true

	fmt.Printf("ES:Disable: %d %s disabled\n", e.ID(), e.Name)
}

func (es *EventSystem) Enable(name string) {
	if es.events[name] == nil {
		log.Fatalf("ES:Disable: Unknow event %s\n", name)
	}

	e := es.events[name]
	hear := es.hears[e.MessageHandlerId]
	if hear == nil {
		return
	}
	hear.Active = false

	fmt.Printf("ES:Disable: %d %s enabled\n", e.ID(), e.Name)
}

func (es *EventSystem) Debug() {
	fmt.Printf("*** Event Manager DEBUG ***\n")
	fmt.Printf("Created: %d\n", len(es.events))
	for name, e := range es.events {
		fmt.Printf("\t- %s: %p, Dispatched: %d, Listened: %t\n", name, e, e.Dispatched, e.Listened)
	}
	fmt.Printf("Heard: %d\n", len(es.hears))
	for _, e := range es.hears {
		fmt.Printf("\t- %p %s\n", e, e.Name)
	}
	fmt.Printf("\n")
}
