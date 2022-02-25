package components

import (
	"math/rand"
	"time"
)

var GameSeed = time.Now().UnixNano()
var GameRand = rand.New(rand.NewSource(GameSeed))
