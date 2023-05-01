package ui

import "github.com/gazercloud/gazerui/opengl/gl11/gl"

func initOpenGL11() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
}
