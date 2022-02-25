package components

import (
	"fmt"
	"github.com/jinzhu/copier"
	"log"
)

type Targets []*GameEntity

type GameEntity struct {
	Name    string
	Type    string
	SubType string
	Sprite  string
	Entity  *Entity
}

func (ge *GameEntity) Copy() *GameEntity {
	cpy := GameEntity{}

	err := copier.Copy(&cpy, ge)
	if err != nil {
		log.Fatalf(fmt.Sprintf("GameEntity:Copy - Failed to copy %s : %s\n", ge.Name, err))
	}

	return &cpy
}
