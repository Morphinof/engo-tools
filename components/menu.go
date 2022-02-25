package components

import "github.com/EngoEngine/engo/common"

type Menu struct {
	common.SpaceComponent
	Name      string
	Cursor    *Entity
	Container *Entity
	Selected  int
	Disabled  EntityArray
	Font      *common.Font
}
