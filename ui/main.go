package ui

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

type Point struct {
	X int
	Y int
}

type MouseCursor int

const (
	MouseCursorNotDefined MouseCursor = 0
	MouseCursorArrow      MouseCursor = 1
	MouseCursorPointer    MouseCursor = 2
	MouseCursorResizeHor  MouseCursor = 3
	MouseCursorResizeVer  MouseCursor = 4
	MouseCursorIBeam      MouseCursor = 5
)

var CursorArrow *glfw.Cursor
var CursorPointer *glfw.Cursor
var CursorResizeHor *glfw.Cursor
var CursorResizeVer *glfw.Cursor
var CursorIBeam *glfw.Cursor

var UseOpenGL33 bool = false
var ServiceDrawBorders = false

func InitUI() error {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	if UseOpenGL33 {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
		glfw.WindowHint(glfw.Samples, 4)
		//InitOpenGL33()
	} else {
		glfw.WindowHint(glfw.ContextVersionMajor, 1)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		initOpenGL11()
	}

	CursorArrow = glfw.CreateStandardCursor(glfw.ArrowCursor)
	CursorPointer = glfw.CreateStandardCursor(glfw.HandCursor)
	CursorResizeHor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	CursorResizeVer = glfw.CreateStandardCursor(glfw.VResizeCursor)
	CursorIBeam = glfw.CreateStandardCursor(glfw.IBeamCursor)

	return nil
}

func InitUISystem() {
	runtime.LockOSThread()

	if err := InitUI(); err != nil {
		fmt.Println(err)
		return
	}

}
