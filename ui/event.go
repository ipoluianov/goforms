package ui

type IEvent interface {
	SetPosX(x int)
	SetPosY(y int)
	PosX() int
	PosY() int
}
