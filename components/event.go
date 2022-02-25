package components

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type Event struct {
	ecs.BasicEntity
	engo.MessageHandlerId
	Name       string
	Data       map[string]any
	Dispatched int
	Listened   bool
	Active     bool
}

func (e Event) Type() string {
	return e.Name
}

type Events map[string]*Event
