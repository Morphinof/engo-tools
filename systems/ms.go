package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image/color"
	"log"
	"tools/components"
)

const (
	MenuItemStartX    = 30
	MenuItemStartY    = 25
	MenuItemTopMargin = 10

	EventMenuItemClicked = "EventMenuItemClicked"
	EventMenuToggle      = "EventMenuToggle"
)

type MenuSystem struct {
	em    *EntityManager
	ev    *EventSystem
	ds    *DragSystem
	ui    *UiSystem
	menus map[string]*components.Menu
}

func (ms *MenuSystem) New(w *ecs.World) {
	ms.menus = make(map[string]*components.Menu)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *EntityManager:
			ms.em = sys
		case *EventSystem:
			ms.ev = sys
		case *DragSystem:
			ms.ds = sys
		case *UiSystem:
			ms.ui = sys
		}
	}

	ms.ev.NewEvent(EventMenuItemClicked)
	ms.ev.NewEvent(EventMenuToggle)

	ms.ev.Listen(EventMenuToggle, func(msg engo.Message) {
		evt := msg.(*components.Event)
		menu := evt.Data["menu"].(*components.Menu)

		// Hide cursor and make children draggable
		for _, m := range ms.menus {
			if m.Container == menu.Container {
				toggle := !m.Container.RenderComponent.Hidden
				m.Container.RenderComponent.Hidden = toggle
				for _, e := range m.Container.Children() {
					ms.em.Get(e).RenderComponent.Hidden = toggle
				}
			}
		}
	})

	ms.ev.Listen(EventStartDrag, func(msg engo.Message) {
		evt := msg.(*components.Event)
		entity := evt.Data["entity"].(*components.Entity)

		// Hide cursor and make children draggable
		for _, m := range ms.menus {
			if m.Container == entity {
				m.Cursor.RenderComponent.Hidden = true
				for _, e := range m.Container.Children() {
					ms.em.Get(e).IsDraggable = true
				}
			}
		}
	})

	ms.ev.Listen(EventStopDrag, func(msg engo.Message) {
		evt := msg.(*components.Event)
		entity := evt.Data["entity"].(*components.Entity)

		// Making sure that parent container is at ZIndex 0 and re-align items in case of speedy drag&drop
		for _, m := range ms.menus {
			if m.Container == entity {
				ms.AlignItems(m)
				m.Container.RenderComponent.StartZIndex = 0
			}
		}
	})
}

func (ms *MenuSystem) Update(dt float32) {
	if engo.Input.Button("F6").JustPressed() {
		ms.Debug()
	}

	for _, m := range ms.menus {
		//if m.Container.MouseComponent.RightClicked {
		//	ms.es.Dispatch(EventMenuToggle, map[string]interface{}{
		//		"menu": m,
		//	})
		//	return
		//}

		for i, e := range m.Container.Children() {
			ce := ms.em.Get(e)

			if ms.IsDisabled(m, ce) {
				continue
			}

			if ce.MouseComponent.Clicked {
				m.Cursor.RenderComponent.Hidden = true
				ms.SelectIndex(m, i)

				ms.ev.Dispatch(EventMenuItemClicked, map[string]any{
					"menu":  m,
					"index": i,
				})
				return
			}

			if ce.MouseComponent.Enter {
				m.Cursor.RenderComponent.Hidden = false
				m.Cursor.Position = ce.Position
				m.Cursor.Position.X -= m.Cursor.Width
				m.Cursor.Position.Y += ce.RenderComponent.Drawable.Height() / 2
				engo.SetCursor(engo.CursorHand)
			}

			if ce.MouseComponent.Leave {
				m.Cursor.RenderComponent.Hidden = true
				engo.SetCursor(engo.CursorNone)
			}
		}
	}
}

func (ms *MenuSystem) Remove(e ecs.BasicEntity) {
	ms.em.Remove(e)
}

func (ms *MenuSystem) RemoveItems(menu *components.Menu) {
	if ms.menus[menu.Name] == nil {
		log.Fatalf("MS:RemoveItems - Unknown menu %s\n", menu.Name)
	}

	for _, e := range menu.Container.Children() {
		menu.Container.RemoveChild(&e)
		ms.em.Remove(e)
	}
}

func (ms *MenuSystem) NewMenu(name string, spaceComponent common.SpaceComponent, container *common.Texture, cursor *common.Texture, items []string, font *common.Font, draggable bool) *components.Menu {
	menu := &components.Menu{
		Name:           name,
		SpaceComponent: spaceComponent,
		Container:      ms.em.NewEntity(),
		Cursor:         ms.em.NewEntity(),
	}
	menu.Font = font
	menu.Container.Refresh = true
	menu.Container.Ref = fmt.Sprintf("menu-%s-container", menu.Name)
	menu.Container.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{
			X: menu.Position.X,
			Y: menu.Position.Y,
		},
		Width:  menu.Width,
		Height: menu.Height,
	}
	menu.Container.RenderComponent = common.RenderComponent{
		Scale: engo.Point{
			X: 1 * (menu.Container.Width / container.Width()),
			Y: 1 * (menu.Container.Height / container.Height()),
		},
		Drawable: container,
		Color:    nil,
	}
	menu.Container.IsDraggable = draggable

	menu.Cursor.Ref = fmt.Sprintf("menu-%s-cursor", menu.Name)
	menu.Cursor.SpaceComponent = common.SpaceComponent{
		//Position: engo.Point{X: menu.Container.Position.X + 20, Y: menu.Container.Position.Y + 20},
		Width:  30,
		Height: 30,
	}
	menu.Cursor.RenderComponent = common.RenderComponent{
		Drawable: cursor,
		Color:    nil,
		Scale: engo.Point{
			X: 1 * (menu.Cursor.Width / cursor.Width()),
			Y: 1 * (menu.Cursor.Height / cursor.Height()),
		},
	}
	menu.Cursor.RenderComponent.Hidden = true
	menu.Container.SetZIndex(components.LayerUiBackground)
	menu.Cursor.SetZIndex(components.LayerUi)

	ms.SetItems(menu, items)
	ms.em.Add(menu.Container, menu.Cursor)
	ms.menus[menu.Name] = menu

	fmt.Printf("MS:NewMenu - Menu %s created\n", menu.Name)

	return menu
}

func (ms *MenuSystem) AlignItems(menu *components.Menu) {
	y := menu.Container.Position.Y + MenuItemStartY
	for _, c := range menu.Container.Children() {
		e := ms.em.Get(c)
		_, h, _ := menu.Font.TextDimensions(e.Drawable.(common.Text).Text)
		e.Position.X = menu.Container.Position.X + MenuItemStartX
		e.Position.Y = y
		e.IsDraggable = false
		y += float32(h)
	}

	fmt.Printf("MS:NewMenu - Menu %s items aligned\n", menu.Name)
}

func (ms *MenuSystem) SetItems(menu *components.Menu, items []string) {
	startPos := menu.Container.Position
	pos := engo.Point{
		X: startPos.X + MenuItemStartX,
		Y: startPos.Y + MenuItemStartY,
	}

	for i, txt := range items {
		t := ms.em.NewEntity()
		t.Ref = fmt.Sprintf("menu-%s-item-%d", menu.Name, i)
		t.RenderComponent.Drawable = common.Text{
			Font: menu.Font,
			Text: txt,
		}
		w, h, _ := menu.Font.TextDimensions(txt)
		t.RenderComponent.SetShader(common.TextHUDShader)
		t.SpaceComponent = common.SpaceComponent{
			Position: pos,
			Width:    float32(w),
			Height:   float32(h),
		}
		t.SetZIndex(components.LayerUi)

		menu.Container.AppendChild(&t.BasicEntity)

		pos.Y += float32(h + MenuItemTopMargin)
	}
}

func (ms *MenuSystem) SelectIndex(menu *components.Menu, index int) {
	menu.Selected = index
	for i, c := range menu.Container.Children() {
		if i == index {
			e := ms.em.Get(c)
			ms.ui.UpdateText(e, "", components.ColorSelected)
		}
	}
}

func (ms *MenuSystem) DisableItems(menu *components.Menu, indexes ...int) {
	for _, index := range indexes {
		for i, c := range menu.Container.Children() {
			if i == index {
				e := ms.em.Get(c)
				//txt := e.RenderComponent.Drawable.(common.Text)
				//font := ms.ui.GetFont(txt.Font.URL, txt.Font.Size, components.ColorDisabled)
				//e.RenderComponent.Drawable = common.Text{
				//	Font: font,
				//	Text: txt.Text,
				//}
				ms.ui.UpdateText(e, "", components.ColorDisabled)
				menu.Disabled = append(menu.Disabled, e)
			}
		}
	}
}

func (ms *MenuSystem) IsDisabled(menu *components.Menu, e *components.Entity) bool {
	for _, me := range menu.Disabled {
		if me == e {
			return true
		}
	}

	return false
}

func (ms *MenuSystem) Show(menu *components.Menu) {
	menu.Container.Hidden = false
	for _, c := range menu.Container.Children() {
		ge := ms.em.Get(c)
		ge.Hidden = false
	}
}

func (ms *MenuSystem) Hide(menu *components.Menu) {
	if menu == nil {
		return
	}

	menu.Container.Hidden = true
	for _, c := range menu.Container.Children() {
		ge := ms.em.Get(c)
		ge.Hidden = true
	}
}

func (ms *MenuSystem) Reset(menu *components.Menu) {
	if menu == nil {
		return
	}

	if len(menu.Container.Children()) > 0 {
		for _, e := range menu.Container.Children() {
			ge := ms.em.Get(e)
			txt := ge.RenderComponent.Drawable.(common.Text)
			font := ms.ui.GetFont(txt.Font.URL, txt.Font.Size, color.Black)
			ge.RenderComponent.Drawable = common.Text{
				Font: font,
				Text: txt.Text,
			}
		}
	}
}

func (ms *MenuSystem) Clean(menu *components.Menu) {
	if menu.Container == nil {
		fmt.Printf("MS:Clean - menu %s container is nil\n", menu.Name)
		return
	}

	if len(menu.Container.Children()) > 0 {
		for _, c := range menu.Container.Children() {
			ms.em.Remove(c)
			menu.Container.RemoveChild(&c)
		}
	}

	menu.Disabled = make(components.EntityArray, 0)
}

func (ms *MenuSystem) Destroy(menu *components.Menu) {
	fmt.Printf("MS:Destroy - Destroying menu %s\n", menu.Name)
	ms.Clean(menu)

	if menu.Container != nil {
		ms.em.Remove(menu.Container.BasicEntity)
		menu.Container = nil
	}

	if menu.Cursor != nil {
		ms.em.Remove(menu.Cursor.BasicEntity)
		menu.Cursor = nil
	}

	delete(ms.menus, menu.Name)
}

func (ms *MenuSystem) Debug() {
	fmt.Printf("*** Menu System DEBUG ***\n")
	fmt.Printf("Instances: %d\n", len(ms.menus))
	for _, m := range ms.menus {
		fmt.Printf("\t- Menu %s\n", m.Name)
		for i, c := range m.Container.Children() {
			e := ms.em.Get(c)
			txt := e.RenderComponent.Drawable.(common.Text)
			fmt.Printf("\t\t- %d %s - Disabled: %t\n", i, txt.Text, ms.IsDisabled(m, e))
		}
	}
	fmt.Printf("\n")
}
