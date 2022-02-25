package main

import (
	"fmt"
	"github.com/EngoEngine/engo"
	"tools/scenes"
)

const (
	Project      = "Tools"
	Version      = 0
	Subversion   = 0
	WindowWidth  = 1280
	WindowHeight = 900
)

func main() {
	opts := engo.RunOptions{
		Title:  fmt.Sprintf("%s v%d.%d", Project, Version, Subversion),
		Width:  WindowWidth,
		Height: WindowHeight,
		//StandardInputs: true,
		//HeadlessMode:   true,
		ScaleOnResize: true,
		NotResizable:  true,
	}

	engo.Run(opts, &scenes.DebugScene{})
}
