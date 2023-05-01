package ui

import (
	"image/color"

	"github.com/golang-collections/collections/stack"
	"github.com/ipoluianov/goforms/opengl/gl11/gl"
)

type ClipSettings struct {
	x      int
	y      int
	width  int
	height int
}

type DrawContextGL struct {
	currentColor color.Color

	fontFamily    string
	fontSize      float64
	fontBold      bool
	fontItalic    bool
	fontUnderline bool

	windowWidth  int
	windowHeight int

	currentClipSettings ClipSettings
	stackClipSettings   stack.Stack
}

func NewDrawContext(windowWidth, windowHeight int) *DrawContextGL {
	var c DrawContextGL
	c.windowWidth = windowWidth
	c.windowHeight = windowHeight
	c.currentClipSettings.x = 0
	c.currentClipSettings.y = 0
	c.currentClipSettings.width = windowWidth
	c.currentClipSettings.height = windowHeight
	c.setViewport()
	return &c
}

func (c *DrawContextGL) SetColor(col color.Color) {
	c.currentColor = col
}

func (c *DrawContextGL) SetFont(fontFamily string, fontSize float64, bold, italic, underline bool) {
	c.fontFamily = fontFamily
	c.fontSize = fontSize
	c.fontBold = bold
	c.fontItalic = italic
	c.fontUnderline = underline
}

func (c *DrawContextGL) Clip(x, y, width, height int) {
	c.stackClipSettings.Push(c.currentClipSettings)
	c.currentClipSettings.x = x
	c.currentClipSettings.y = y
	c.currentClipSettings.width = width
	c.currentClipSettings.height = height
	c.setViewport()
}

func (c *DrawContextGL) setViewport() {
	gl.Viewport(int32(c.currentClipSettings.x), int32(c.windowHeight-c.currentClipSettings.y-c.currentClipSettings.height), int32(c.currentClipSettings.width), int32(c.currentClipSettings.height))
}

func (c *DrawContextGL) UnClip() {
	c.currentClipSettings = c.stackClipSettings.Peek().(ClipSettings)
	c.setViewport()
}

func (c *DrawContextGL) xyTo01(x, y int) (float32, float32) {
	dx := 2 / float64(c.currentClipSettings.width)
	dy := 2 / float64(c.currentClipSettings.height)
	x1d := float32(dx*float64(x)) - 1
	y1d := float32(dy*float64(y)) - 1
	y1d = -y1d // OpenGL coordinates are inverted by Y
	return x1d, y1d
}

/*func (c *DrawContextGL) FillRect(x, y, width, height int) {
	temp.checkProgram()

	x1d, y1d := c.xyTo01(x, y)
	x2d, y2d := c.xyTo01(width+x, y)
	x3d, y3d := c.xyTo01(width+x, height+y)
	x4d, y4d := c.xyTo01(x, height+y)

	vertices := []float32{
		x1d, y1d, 0.0,
		x2d, y2d, 0.0,
		x3d, y3d, 0.0,
		x3d, y3d, 0.0,
		x4d, y4d, 0.0,
		x1d, y1d, 0.0,
	}

	r, g, b, a := c.currentColor.RGBA()
	gl.UseProgram(temp.SimpleProgram)
	gl.Uniform4f(gl.GetUniformLocation(temp.SimpleProgram, gl.Str("pixelColor\x00")), float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535)

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// specify the format of our vertex input
	// (shader) input 0
	// vertex has size 3
	// vertex items are of type FLOAT
	// do not normalize (already done)
	// stride of 3 * sizeof(float) (separation of vertices)
	// offset of where the position data starts (0 for the beginning)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	gl.BindVertexArray(VAO)           // bind data
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // perform draw call
	gl.BindVertexArray(0)             // unbind data (so we don't mistakenly use/modify it)

	gl.DeleteBuffers(1, &VBO)
	gl.DeleteVertexArrays(1, &VAO)
}

func (c *DrawContextGL) DrawRect(x, y, width, height int, lineWidth int) {

}*/

/*func (c *DrawContextGL) DrawLine(x1, y1, x2, y2 int, width int) {
	temp.checkProgram()

	vertices := c.calcLineRect(float32(x1), float32(y1), float32(x2), float32(y2), width)
	vertices[0], vertices[1] = c.xyTo01(int(vertices[0]), int(vertices[1]))
	vertices[3], vertices[4] = c.xyTo01(int(vertices[3]), int(vertices[4]))
	vertices[6], vertices[7] = c.xyTo01(int(vertices[6]), int(vertices[7]))
	vertices[9], vertices[10] = c.xyTo01(int(vertices[9]), int(vertices[10]))
	vertices[12], vertices[13] = c.xyTo01(int(vertices[12]), int(vertices[13]))
	vertices[15], vertices[16] = c.xyTo01(int(vertices[15]), int(vertices[16]))

	//x1d, y1d := c.xyTo01(x1, y1)
	//x2d, y2d := c.xyTo01(x2, y2)

	//x1d, y1d := float32(-0.5), float32(-0.5)
	//x2d, y2d := float32(0.5), float32(0.5)

	r, g, b, a := c.currentColor.RGBA()
	gl.UseProgram(temp.SimpleProgram)
	gl.Uniform4f(gl.GetUniformLocation(temp.SimpleProgram, gl.Str("pixelColor\x00")), float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535)

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// specify the format of our vertex input
	// (shader) input 0
	// vertex has size 3
	// vertex items are of type FLOAT
	// do not normalize (already done)
	// stride of 3 * sizeof(float) (separation of vertices)
	// offset of where the position data starts (0 for the beginning)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	gl.BindVertexArray(VAO)           // bind data
	gl.DrawArrays(gl.TRIANGLES, 0, 6) // perform draw call
	gl.BindVertexArray(0)             // unbind data (so we don't mistakenly use/modify it)

	gl.DeleteBuffers(1, &VBO)
	gl.DeleteVertexArrays(1, &VAO)
}*/

func (c *DrawContextGL) DrawEllipse(x, y, width, height int) {
}

func (c *DrawContextGL) DrawPolygon(x, y, width, height int) {
}

func (c *DrawContextGL) DrawText(x, y, width, height int, text string) {
}
