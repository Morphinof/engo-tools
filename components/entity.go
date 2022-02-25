package components

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/jinzhu/copier"
	"log"
)

type EntityArray []*Entity
type EntityMap map[uint64]*Entity

type Chained struct {
	Next *Entity
	Prev *Entity
}

type Draggable struct {
	XOff        float32
	YOff        float32
	Drag        bool
	IsDraggable bool
}

type Entity struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.MouseComponent
	common.RenderComponent
	Draggable
	Chained
	Ref     string
	Refresh bool
}

func (e *Entity) StartDrag() {
	e.XOff = engo.Input.Mouse.X - e.SpaceComponent.Position.X
	e.YOff = engo.Input.Mouse.Y - e.SpaceComponent.Position.Y
	e.Drag = true
}

func (e *Entity) StopDrag() {
	e.Drag = false
	e.XOff = 0
	e.YOff = 0
}

func (e *Entity) Copy() *Entity {
	cpy := Entity{}

	err := copier.Copy(&cpy, e)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Entity:Copy - Failed to copy %d : %s\n", e.ID(), err))
	}

	return &cpy
}

func (e *Entity) RemoveFromChain() {
	prev := e.Prev
	next := e.Next
	if prev != nil {
		prev.Next = next
		if next != nil {
			next.Prev = prev
		}
	}
}
