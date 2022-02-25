package scenes

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"image/color"
	"tools/components"
	"tools/systems"
)

var GreenLight = color.RGBA{R: 178, G: 192, B: 168, A: 255}
var Green = color.RGBA{R: 118, G: 154, B: 103, A: 255}
var GreenDark = color.RGBA{R: 52, G: 93, B: 81, A: 255}

const (
	EventClicked = "EventClicked"
)

type DebugScene struct {
	world *ecs.World
	es    systems.EventSystem
	em    systems.EntityManager
	ds    systems.DragSystem
	ui    systems.UiSystem
	ms    systems.MenuSystem
}

// Preload initializes assets
func (ds *DebugScene) Preload() {
	fmt.Printf("Debug scene preload\n")
}

// Setup function
func (ds *DebugScene) Setup(u engo.Updater) {
	fmt.Printf("Debug scene setup\n")

	common.SetBackground(GreenLight)
	ds.setupKeys()
	ds.ui.LoadFonts()
	ds.ui.LoadAssets()

	ds.world = u.(*ecs.World)
	ds.world.AddSystem(&ds.es)
	ds.world.AddSystem(&ds.em)
	ds.world.AddSystem(&ds.ds)
	ds.world.AddSystem(&ds.ui)
	ds.world.AddSystem(&ds.ms)

	// Listen the EventMenuItemClicked event
	ds.es.Listen(systems.EventMenuItemClicked, func(m engo.Message) {
		evt := m.(*components.Event)
		menu := evt.Data["menu"].(*components.Menu)
		index := evt.Data["index"].(int)

		fmt.Printf("Menu %s item %d clicked\n", menu.Name, index)
	})

	// Register the EventClicked event
	ds.es.NewEvent(EventClicked)

	// Listen the EventClicked event
	ds.es.Listen(EventClicked, func(m engo.Message) {
		evt := m.(*components.Event)
		entity := evt.Data["entity"].(*components.Entity)

		fmt.Printf("Target %d clicked\n", entity.ID())
	})

	ds.tests()
}

func (ds *DebugScene) setupKeys() {
	mapping := map[string]engo.Key{
		"Q":   engo.KeyQ,
		"E":   engo.KeyE,
		"W":   engo.KeyW,
		"TAB": engo.KeyTab,
		"F1":  engo.KeyF1,
		"F2":  engo.KeyF2,
		"F3":  engo.KeyF3,
		"F4":  engo.KeyF4,
		"F5":  engo.KeyF5,
		"F6":  engo.KeyF6,
		"F7":  engo.KeyF7,
		"F8":  engo.KeyF8,
		"F9":  engo.KeyF9,
		"F10": engo.KeyF10,
		"F11": engo.KeyF11,
		"F12": engo.KeyF12,
	}

	input := engo.Input
	for key, engoKey := range mapping {
		input.RegisterButton(key, engoKey)
	}
}

// Type returns the type
func (ds *DebugScene) Type() string {
	return "DebugScene"
}

func (ds *DebugScene) tests() {
	ds.TestMenuComponent()
	ds.TestMenuRefreshItems()
}

func (ds *DebugScene) TestMenuComponent() {
	container := ds.ui.LoadSprite("box")
	cursor := ds.ui.LoadSprite("cursor")
	items := []string{
		"Attack",
		"Use",
		"Cast",
	}
	font := ds.ui.GetFont("Roboto-Regular.ttf", 26, color.Black)

	width, height := float32(250.0), float32(350.0)
	menu := ds.ms.NewMenu("test-menu", common.SpaceComponent{
		Position: engo.Point{
			X: engo.WindowWidth()/2 - width/2,
			Y: engo.WindowHeight()/2 - height/2,
		},
		Width:  width,
		Height: height,
	}, container, cursor, items, font, false)
	menu.Container.RenderComponent.SetZIndex(components.LayerUi)
}

func (ds *DebugScene) TestMenuRefreshItems() {
	container := ds.ui.LoadSprite("box")
	cursor := ds.ui.LoadSprite("cursor")
	items := []string{
		"Attack",
		"Use",
		"Cast",
	}
	font := ds.ui.GetFont("8-bit-hud.ttf", 30, color.Black)
	width, height := float32(250.0), float32(350.0)
	menu := ds.ms.NewMenu("test-menu-refresh-items", common.SpaceComponent{
		Position: engo.Point{
			X: 50,
			Y: 50,
		},
		Width:  width,
		Height: height,
	}, container, cursor, items, font, false)

	newItems := []string{
		"Claw",
		"Fire",
		"Recover",
	}

	ds.ms.RemoveItems(menu)
	ds.ms.SetItems(menu, newItems)
	menu.Container.RenderComponent.SetZIndex(components.LayerUi)
}
