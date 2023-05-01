package ui

import "github.com/ipoluianov/goforms/opengl/gl11/gl"

func initOpenGL11() {
	err := gl.Init()
	if err != nil {
		panic(err)
	}
}
