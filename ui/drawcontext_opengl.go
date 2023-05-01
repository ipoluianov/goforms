package ui

import (
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ipoluianov/goforms/canvas"
)

func init() {
	texturesMap = make(map[int]*image.RGBA)
}

type DrawContextOpenGL struct {
	DrawContextBase
}

func NewDrawContextOpenGL(window *glfw.Window) *DrawContextOpenGL {
	var c DrawContextOpenGL
	c.Window = window
	c.WindowWidth, c.WindowHeight = window.GetSize()

	c.DrawContextBase.InitBase()

	c.setViewport()
	return &c
}

var vao uint32

func (c *DrawContextOpenGL) Init() {
	c.DrawContextBase.InitBase()
	c.setViewport()
	gl.ClearColor(0, 0.3, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	checkProgram()
	gl.UseProgram(SimpleProgram)
}

func (c *DrawContextOpenGL) GG() *gg.Context {
	return nil
}

func (c *DrawContextOpenGL) GraphContextImage() *image.RGBA {
	return nil
}

func (c *DrawContextOpenGL) State() canvas.CanvasDirectState {
	return canvas.CanvasDirectState{}
}

func (c *DrawContextOpenGL) setViewport() {
	gl.Viewport(int32(c.CurrentClipSettings.x), int32(c.WindowHeight-c.CurrentClipSettings.y-c.CurrentClipSettings.height), int32(c.CurrentClipSettings.width), int32(c.CurrentClipSettings.height))
}

func (c *DrawContextOpenGL) xyTo01(x, y int) (float32, float32) {
	dx := 2 / float64(c.CurrentClipSettings.width)
	dy := 2 / float64(c.CurrentClipSettings.height)
	x1d := float32(dx*float64(x)) - 1
	y1d := float32(dy*float64(y)) - 1
	y1d = -y1d // OpenGL coordinates are inverted by Y
	return x1d, y1d
}

func (c *DrawContextOpenGL) DrawRect(x, y, width, height int) {
	if width < 0 {
		x = x + width
		width = -width
	}

	if height < 0 {
		y = y + height
		height = -height
	}

	x1 := x
	y1 := y
	x2 := x + width
	y2 := y
	x3 := x + width
	y3 := y + height
	x4 := x
	y4 := y + height

	lineWidth := c.StrokeWidth
	c.StrokeWidth = 5

	c.DrawLine(x1-lineWidth/2, y1, x2+lineWidth/2, y2)
	c.DrawLine(x2, y2, x3, y3)
	c.DrawLine(x3+lineWidth/2, y3, x4-lineWidth/2, y4)
	c.DrawLine(x4, y4, x1, y1)
}

func (c *DrawContextOpenGL) FillRect(x, y, width, height int) {
	checkProgram()

	x1, y1 := c.xyTo01(x, y)
	x2, y2 := c.xyTo01(x+width, y)
	x3, y3 := c.xyTo01(x+width, y+height)
	x4, y4 := c.xyTo01(x, y+height)

	vertices := []float32{
		x1, y1, 0,
		x2, y2, 0,
		x4, y4, 0,
		x2, y2, 0,
		x3, y3, 0,
		x4, y4, 0,
	}

	r, g, b, a := c.CurrentColor.RGBA()
	gl.UseProgram(SimpleProgram)
	r = 65535
	gl.Uniform4f(gl.GetUniformLocation(SimpleProgram, gl.Str("pixelColor\x00")), float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/3))

	gl.DeleteBuffers(1, &vao)
	gl.DeleteVertexArrays(1, &vao)
}

func (c *DrawContextOpenGL) DrawImage(x, y, width, height int, img image.Image) {
}

func (c *DrawContextOpenGL) DrawText(x, y, width, height int, text string) {
}

func (c *DrawContextOpenGL) TranslatedX() int {
	return c.State().TranslateX
}

func (c *DrawContextOpenGL) TranslatedY() int {
	return c.State().TranslateY
}

func (c *DrawContextOpenGL) ClippedRegion() (int, int, int, int) {
	return 0, 0, 0, 0
}

func (c *DrawContextOpenGL) DrawLine(x1, y1, x2, y2 int) {

	checkProgram()

	vertices := c.calcLineRect(float32(x1), float32(y1), float32(x2), float32(y2), float32(c.StrokeWidth))
	vertices[0], vertices[1] = c.xyTo01(int(vertices[0]), int(vertices[1]))
	vertices[3], vertices[4] = c.xyTo01(int(vertices[3]), int(vertices[4]))
	vertices[6], vertices[7] = c.xyTo01(int(vertices[6]), int(vertices[7]))
	vertices[9], vertices[10] = c.xyTo01(int(vertices[9]), int(vertices[10]))
	vertices[12], vertices[13] = c.xyTo01(int(vertices[12]), int(vertices[13]))
	vertices[15], vertices[16] = c.xyTo01(int(vertices[15]), int(vertices[16]))

	r, g, b, a := c.CurrentColor.RGBA()
	gl.UseProgram(SimpleProgram)
	g = 65535
	gl.Uniform4f(gl.GetUniformLocation(SimpleProgram, gl.Str("pixelColor\x00")), float32(r)/65535, float32(g)/65535, float32(b)/65535, float32(a)/65535)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/3))

	gl.DeleteBuffers(1, &vao)
	gl.DeleteVertexArrays(1, &vao)
}

func (c *DrawContextOpenGL) DrawEllipse(x, y, width, height int) {
	checkProgram()

	numSegments := 10
	r := float32(0.1)
	cx, cy := c.xyTo01(x, y)

	ll := make([]float32, 0)
	for ii := 0; ii < numSegments; ii++ {
		theta := 2.0 * 3.1415926 * float64(ii) / float64(numSegments)
		x := r * float32(math.Cos(theta))
		y := r * float32(math.Sin(theta))
		ll = append(ll, x+cx)
		ll = append(ll, y+cy)
		ll = append(ll, 0)
	}

	gl.UseProgram(SimpleProgram)
	//gl.Uniform4f(gl.GetUniformLocation(SimpleProgram, gl.Str("pixelColor\x00")), 1, 0, 0, 1)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(ll), gl.Ptr(ll), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(ll)/3))

	gl.DeleteBuffers(1, &vbo)
	gl.DeleteVertexArrays(1, &vao)

}

func (c *DrawContextOpenGL) calcLineRect(x1, y1, x2, y2 float32, width float32) []float32 {
	//width -= 1
	//kxy := float32(c.currentClipSettings.width) / float32(c.currentClipSettings.height)
	width /= 2

	// kxy = 1 / kxy

	//kxy = kxy

	deltaX := x2 - x1
	deltaY := y2 - y1

	tnA := deltaY / deltaX

	katY := float32(width) * float32(math.Cos(math.Atan(float64(tnA))))
	katX := float32(math.Sqrt(math.Abs(float64(float64(width*width) - float64(katY*katY)))))

	//katY = float32((2 / float64(c.currentClipSettings.height)) * float64(katY)) // To OpenGL coordinates
	//katX = float32((2 / float64(c.currentClipSettings.width)) * float64(katX)) // To OpenGL coordinates

	/*fmt.Println("kxy:", kxy , " katY:", katY, " katX:", katX)

	katX = katX / kxy
	katY = katY * kxy*/

	if deltaX <= 0 && deltaY <= 0 {
		katY = -katY
	}

	if deltaX <= 0 && deltaY > 0 {
		katX = -katX
		katY = -katY
	}

	if deltaX > 0 && deltaY > 0 {
		katX = -katX
	}

	x1d := x1 - katX
	y1d := y1 - katY
	x2d := x1 + katX
	y2d := y1 + katY

	x3d := x2 + katX
	y3d := y2 + katY
	x4d := x2 - katX
	y4d := y2 - katY

	vertices := []float32{
		x1d, y1d, 0.0,
		x2d, y2d, 0.0,
		x3d, y3d, 0.0,

		x3d, y3d, 0.0,
		x4d, y4d, 0.0,
		x1d, y1d, 0.0,
	}
	return vertices
}

func (c *DrawContextOpenGL) MeasureText(text string) (int, int) {
	return 0, 0
}

func (c *DrawContextOpenGL) Save() {
	c.StackClipSettings.Push(c.CurrentClipSettings)
}

func (c *DrawContextOpenGL) Translate(x, y int) {
	c.CurrentClipSettings.x += x
	c.CurrentClipSettings.y += y
	c.setViewport()
}

func (c *DrawContextOpenGL) Clip(x, y, width, height int) {
	c.CurrentClipSettings.x += x
	c.CurrentClipSettings.y += y
	c.CurrentClipSettings.width = width
	c.CurrentClipSettings.height = height
	c.setViewport()
}

func (c *DrawContextOpenGL) ClipIn(x, y, width, height int) {

}

func (c *DrawContextOpenGL) Load() {
	c.CurrentClipSettings = c.StackClipSettings.Peek().(ClipSettings)
	c.StackClipSettings.Pop()
	c.setViewport()
}

func (c *DrawContextOpenGL) Finish() {
}

var SimpleProgramLoaded = false
var SimpleProgram uint32 = 0

func checkProgram() {
	if !SimpleProgramLoaded {
		vertexShader, err := compileShader(SimpleProgramVertex, gl.VERTEX_SHADER)
		if err != nil {
			return
		}

		fragmentShader, err := compileShader(SimpleProgramFragment, gl.FRAGMENT_SHADER)
		if err != nil {
			return
		}

		program := gl.CreateProgram()

		gl.AttachShader(program, vertexShader)
		gl.AttachShader(program, fragmentShader)
		//gl.AttachShader(program, geomShader)
		gl.LinkProgram(program)

		var status int32
		gl.GetProgramiv(program, gl.LINK_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

			panic("Error while compiling shaders")

			return
		}

		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)
		//gl.DeleteShader(geomShader)

		SimpleProgram = program
		SimpleProgramLoaded = true

		return

	}
}

// compileShader compiles the shader program
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var SimpleProgramFragment = `#version 330 core

out vec4 color;

uniform vec4 pixelColor;

void main()
{
    color = pixelColor;
}` + "\x00"

var SimpleProgramVertex = `#version 330 core

layout (location = 0) in vec3 position;

void main()
{
    gl_Position = vec4(position.x, position.y, position.z, 1.0);
}` + "\x00"
