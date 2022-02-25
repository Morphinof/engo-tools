package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"log"
	"tools/components"
)

type EntityManager struct {
	world        *ecs.World
	renderSystem *common.RenderSystem
	mouseSystem  *common.MouseSystem
	instances    components.EntityMap
	sent         []uint64
	buffer       components.EntityArray
	flushed      int
}

func (em *EntityManager) New(w *ecs.World) {
	em.world = w
	em.renderSystem = new(common.RenderSystem)
	em.mouseSystem = new(common.MouseSystem)
	em.world.AddSystem(em.renderSystem)
	em.world.AddSystem(em.mouseSystem)
	em.instances = make(components.EntityMap)
	em.sent = make([]uint64, 0)
}

func (em *EntityManager) Update(dt float32) {
	if engo.Input.Button("F1").JustPressed() {
		em.Debug()
	}

	l := len(em.buffer)
	if l > 0 {
		em.Flush()
	}

	for _, e := range em.instances {
		if e.Refresh {
			em.Refresh(e)
		}
	}
}

func (em *EntityManager) Remove(e ecs.BasicEntity) {
	instance := em.instances[e.ID()]
	if instance == nil {
		fmt.Printf("EM:Remove - Failed, unknown entity %d\n", e.ID())
		return
	}

	instance.RemoveFromChain()

	count := len(e.Children())
	if count > 0 {
		//fmt.Printf("EM:Remove - Entity %d remove %d children\n", e.ID(), count)
		for _, c := range e.Children() {
			em.Remove(c)
		}
	}

	em.renderSystem.Remove(e)
	for i, _id := range em.sent {
		if _id == e.ID() {
			em.sent = append(em.sent[:i], em.sent[i+1:]...)
		}
	}
	delete(em.instances, e.ID())

	//fmt.Printf("EM:Remove - Entity %d removed\n", e.ID())
}

func (em *EntityManager) Add(entities ...*components.Entity) {
	var prev *components.Entity
	for _, e := range entities {
		em.buffer = append(em.buffer, e)
		em.instances[e.ID()] = e
		count := len(e.Children())
		if count > 0 {
			prev = e
			for _, c := range e.Children() {
				ce := em.instances[c.ID()]
				ce.Prev = prev
				prev.Next = ce
				prev = ce
			}
		}
	}
}

func (em *EntityManager) Copy(entity *components.Entity) *components.Entity {
	cpy := entity.Copy()
	cpy.BasicEntity = ecs.NewBasic()
	cpy.Ref = fmt.Sprintf("copy-%s", entity.Ref)
	em.Add(cpy)
	return cpy
}

func (em *EntityManager) Get(entities ...ecs.BasicEntity) *components.Entity {
	for _, e := range entities {
		if em.instances[e.ID()] != nil {
			return em.instances[e.ID()]
		}
	}

	return nil
}

func (em *EntityManager) GetInstances() components.EntityMap {
	return em.instances
}

func (em *EntityManager) Flush() {
	for _, e := range em.buffer {
		if e.Ref == "" {
			fmt.Printf("EM:Flush - Warning ! Entity %d has no reference\n", e.ID())
		}

		if em.instances[e.ID()] == nil {
			fmt.Printf("EM:Flush - skip entity %d (%s), removed or not managed\n", e.ID(), e.Ref)
			if e.Parent() != nil {
				parent := em.instances[e.Parent().ID()]
				fmt.Printf("EM:Flush - skipped entity parent is: %s\n", parent.Ref)
			} else {
				fmt.Printf("EM:Flush - skipped entity has no parent\n")
			}
			continue
		}

		em.add(e)

		count := len(e.Children())
		if count > 0 {
			for _, c := range e.Children() {
				ce := em.instances[c.ID()]
				if em.instances[e.ID()] == nil {
					fmt.Printf("EM:Flush - skip child %d, removed or not managed\n", ce.ID())
					continue
				}

				em.add(ce)
			}
		}
	}

	em.buffer = make(components.EntityArray, 0)
	fmt.Printf("EM:Update - %d entities flushed\n", em.flushed)
	em.flushed = 0
}

func (em *EntityManager) Refresh(e *components.Entity) {
	count := len(e.Children())
	if count > 0 {
		for _, c := range e.Children() {
			if em.renderSystem.EntityExists(&c) == -1 {
				ce := em.Get(c)
				if ce == nil {
					fmt.Printf("EM:Refresh - skip child %d of entity %d, removed or not managed\n", c.ID(), e.ID())
					continue
				}

				//fmt.Printf("EM:Refresh - new unmanaged child %d of entity %d detected\n", ce.ID(), e.ID())
				em.add(ce)
			}
		}
	}
}

func (em *EntityManager) Display(e *components.Entity, hidden bool) {
	e.Hidden = hidden
	count := len(e.Children())
	if count > 0 {
		for _, c := range e.Children() {
			ge := em.Get(c)
			ge.Hidden = hidden
		}
	}
}

func (em *EntityManager) add(e *components.Entity) {
	if &e.MouseComponent == nil {
		log.Fatalf("EM:add - Entity %d missing MouseComponent\n", e.ID())
	}

	if &e.SpaceComponent == nil {
		log.Fatalf("EM:add - Entity %d missing SpaceComponent\n", e.ID())
	}

	if &e.RenderComponent == nil {
		log.Fatalf("EM:add - Entity %d missing RenderComponent\n", e.ID())
	}

	em.mouseSystem.Add(&e.BasicEntity, &e.MouseComponent, &e.SpaceComponent, &e.RenderComponent)
	em.renderSystem.Add(&e.BasicEntity, &e.RenderComponent, &e.SpaceComponent)
	em.sent = append(em.sent, e.ID())
	em.flushed++
}

func (em *EntityManager) NewEntity() *components.Entity {
	e := &components.Entity{
		BasicEntity:     ecs.NewBasic(),
		MouseComponent:  common.MouseComponent{},
		RenderComponent: common.RenderComponent{},
		SpaceComponent:  common.SpaceComponent{},
	}

	em.instances[e.ID()] = e

	return e
}

func (em *EntityManager) Debug() {
	fmt.Printf("*** Entity Manager DEBUG ***\n")
	fmt.Printf("Instances: %d\n", len(em.instances))
	for id, e := range em.instances {
		fmt.Printf("\t- %d: %s %p \n", id, e.Ref, e)
	}
	fmt.Printf("Sent: %d\n", len(em.sent))
	sent := ""
	for i, id := range em.sent {
		if i > 0 {
			sent += ", "
		}
		sent += fmt.Sprintf("%d", id)
	}
	fmt.Printf("[%s]\n", sent)
	fmt.Printf("Buffer: %d\n", len(em.buffer))
	for _, e := range em.buffer {
		fmt.Printf("\t- %d: %p %s\n", e.ID(), e, e.Ref)
	}
	fmt.Printf("Chains\n")
	for _, e := range em.instances {
		if len(e.Children()) > 0 {
			str := fmt.Sprintf("%d", e.ID())
			for _, c := range e.Children() {
				str += fmt.Sprintf(" -> %d", c.ID())
			}
			fmt.Printf("\t - %s\n", str)
		}
	}
	fmt.Printf("\n")
}
