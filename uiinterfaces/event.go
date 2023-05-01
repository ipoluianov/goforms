package uiinterfaces

type Event interface {
	SetPosX(x int)
	SetPosY(y int)
	PosX() int
	PosY() int
}
