package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

const (
	EventStartDrag = "EventStartDrag"
	EventStopDrag  = "EventStopDrag"
)

type DragSystem struct {
	em *EntityManager
	ev *EventSystem
}

func (ds *DragSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *EntityManager:
			ds.em = sys
		case *EventSystem:
			ds.ev = sys
		}
	}

	ds.ev.NewEvent(EventStartDrag)
	ds.ev.NewEvent(EventStopDrag)
}

func (ds *DragSystem) Update(dt float32) {
	eStartDrag := ds.ev.Get(EventStartDrag)
	eStopDrag := ds.ev.Get(EventStopDrag)
	if eStartDrag == nil || !eStartDrag.Active || eStopDrag == nil || !eStopDrag.Active {
		return
	}

	for _, e := range ds.em.GetInstances() {
		if !e.IsDraggable {
			continue
		}

		if e.MouseComponent.Clicked {
			for _, c := range e.Children() {
				ce := ds.em.Get(c)

				// Prevent drag if a children has been clicked
				if ce.MouseComponent.Clicked {
					fmt.Printf("Cancel drag because %d has been clicked\n", ce.ID())
					return
				}

				if ce != nil {
					ce.StartDrag()
				}
			}
			e.StartDrag()
			ds.ev.Dispatch(EventStartDrag, map[string]any{
				"entity": e,
			})
			return
		}

		if e.MouseComponent.Released && e.Drag {
			for _, c := range e.Children() {
				ce := ds.em.Get(c)
				ce.StopDrag()
			}
			e.StopDrag()
			ds.ev.Dispatch(EventStopDrag, map[string]any{
				"entity": e,
			})
			return
		}

		if e.Drag {
			e.SpaceComponent.Position.Set(engo.Input.Mouse.X-e.XOff, engo.Input.Mouse.Y-e.YOff)
		}
	}
}

func (ds *DragSystem) Disable() {
	ds.ev.Disable(EventStartDrag)
	ds.ev.Disable(EventStopDrag)
	fmt.Printf("DS:Disable - Drag system disabled\n")
}

func (ds *DragSystem) Enable() {
	ds.ev.Enable(EventStartDrag)
	ds.ev.Enable(EventStopDrag)
	fmt.Printf("DS:Disable - Drag system enabled\n")
}

func (ds *DragSystem) Remove(e ecs.BasicEntity) {
	ds.em.Remove(e)
}
