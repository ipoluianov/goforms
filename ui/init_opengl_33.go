package ui

import "github.com/go-gl/gl/v3.3-core/gl"

func InitOpenGL33() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
}
