package components

type Dices struct {
	X     int
	Faces int
}

func (d *Dices) Roll() float64 {
	roll := 1
	rnd := GameRand

	for i := 0; i < d.X; i++ {
		roll += rnd.Intn(d.Faces)
	}

	return float64(roll)
}

func NewDices(faces int, x int) *Dices {
	return &Dices{
		X:     x,
		Faces: faces,
	}
}
