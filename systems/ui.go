package systems

import (
	"bytes"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"io/ioutil"
	"log"
	"tools/components"
)

type UiSystem struct {
	em *EntityManager
}

func (ui *UiSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *EntityManager:
			ui.em = sys
		}
	}
}

func (ui *UiSystem) Update(dt float32) {
}

func (ui *UiSystem) Remove(e ecs.BasicEntity) {
	ui.em.Remove(e)
}

func (ui *UiSystem) LoadFile(file string) {
	err := engo.Files.Load(file)
	if err != nil {
		log.Fatalf("DS:LoadFile - Failed to load file with URL: %v\n", file)
	}
}

func (ui *UiSystem) LoadFonts() {
	files := []string{
		"CN.ttf",
		"Roboto-Regular.ttf",
		"8-bit-pusab.ttf",
		"8-bit-madness.ttf",
		"8-bit-hud.ttf",
		"8-bit-16.ttf",
		"PixelOperatorSC.ttf",
	}
	for _, f := range files {
		err := engo.Files.LoadReaderData(f, bytes.NewReader(goregular.TTF))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("DS:LoadFonts - Font %s loaded\n", f)
	}
}

func (ui *UiSystem) LoadAssets() {
	files, err := ioutil.ReadDir("assets")
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		ui.LoadFile(f.Name())
		fmt.Printf("DS:LoadAssets - Asset %s loaded\n", f.Name())
	}
}

func (ui *UiSystem) LoadSprite(sprite string) *common.Texture {
	texture, err := common.LoadedSprite(sprite + ".png")
	if err != nil {
		log.Fatalf(fmt.Sprintf("DS:LoadSprite: Failed to load \"%s\", %s\n", sprite, err))
		return nil
	}

	return texture
}

func (ui *UiSystem) GetFont(font string, size float64, color color.Color) *common.Font {
	fnt := common.Font{
		URL:  font,
		FG:   color,
		Size: size,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		log.Fatal(fmt.Sprintf("DS:GetFont - %s\n", err.Error()))
		return nil
	}

	return &fnt
}

func (ui *UiSystem) NewText(txt string, sc common.SpaceComponent, size float64, font string, color color.Color) *components.Entity {
	t := ui.em.NewEntity()
	t.Ref = fmt.Sprintf("txt-%d", t.ID())
	t.RenderComponent.Drawable = common.Text{
		Font: ui.GetFont(font, size, color),
		Text: txt,
	}
	t.RenderComponent.SetShader(common.TextHUDShader)
	t.SpaceComponent = sc

	ui.em.Add(t)

	return t
}

func (ui *UiSystem) UpdateText(e *components.Entity, text string, color color.Color) {
	txt := e.RenderComponent.Drawable.(common.Text)
	if text == "" {
		text = txt.Text
	}
	font := ui.GetFont(txt.Font.URL, txt.Font.Size, color)
	e.RenderComponent.Drawable = common.Text{
		Font: font,
		Text: text,
	}
}
