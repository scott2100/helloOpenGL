package main

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"math"
	"runtime"
	"strings"
)

var triangle2 = []float32 {
	//first triangle
	//positions      //colors
	0, 0.5, 0,       1.0, 0.0, 0.0,
	-0.5, -0.5, 0,   0.0, 1.0, 0,0,
	0.5, -0.5, 0,    0.0, 0.0, 1.0,
}

var triangle1 = []float32 {
	//first triangle
	//positions
	0, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
}

/*
var triangle2 = []float32 {
	//second triangle
	1, 1, 0,
	.5, .5 , 0,
	.5, 0, 0,
}

 */

const (

	vertexShaderSourceRed = `
		#version 330 core
		layout (location = 0) in vec3 aPos;
		void main()
		{
			gl_Position = vec4(aPos, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 330 core
		out vec4 FragColor;
		uniform vec4 ourColor;
		void main()
		{
			FragColor = ourColor;
		} 
	` + "\x00"

	vertexShaderSourceMulti = `
		#version 330 core
		layout (location = 0) in vec3 aPos;
		layout (location = 1) in vec3 aColor;

		out vec3 ourColor;

		void main()
		{
			gl_Position = vec4(aPos, 1.0);
			ourColor = aColor;
		}
	` + "\x00"

	fragmentShaderSourceMulti= `
		#version 330 core
		out vec4 FragColor;
		in vec3 ourColor;
		void main()
		{
			FragColor = vec4(ourColor, 1.0);
		} 
	` + "\x00"


)

func main()  {
	runtime.LockOSThread()

	glfw.Init()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)

	window, err := glfw.CreateWindow(640, 480, "Scott Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	vao1 := makeVao(triangle1)
	//vao2 := makeVao(triangle2)
	var vbo2 uint32
	var vao2 uint32
	gl.GenVertexArrays(1, &vao2)
	gl.GenBuffers(1, &vbo2)
	gl.BindVertexArray(vao2)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle2), gl.Ptr(triangle2), gl.STATIC_DRAW)



	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6 * 4, gl.PtrOffset(0))
	//gl.EnableVertexAttribArray(0)
	//gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6 * 4, gl.PtrOffset(3 * 4))
	//gl.EnableVertexAttribArray(1)

	vertexShaderMulti, err := compileShader(vertexShaderSourceMulti, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	vertexShaderRed, err := compileShader(vertexShaderSourceRed, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShaderMulti, err := compileShader(fragmentShaderSourceMulti, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShaderRed)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)


	prog2 := gl.CreateProgram()
	gl.AttachShader(prog2, vertexShaderMulti)
	gl.AttachShader(prog2, fragmentShaderMulti)
	gl.LinkProgram(prog2)


	//render loop
	for !window.ShouldClose() {
		//process commands
		processInput(window)

		//render commands
		gl.ClearColor(.2, .3, .3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(prog2)
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		gl.UseProgram(prog)
		timeValue := glfw.GetTime()
		greenValue := (math.Sin(timeValue) / 2) + 0.5
		vertexColorLocation := gl.GetUniformLocation(prog, gl.Str("ourColor\x00"))
		gl.Uniform4f(vertexColorLocation, 0.0, float32(greenValue), 0.0, 1.0)
		gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle1)/3))


		// check and call events and swap the buffers
		glfw.PollEvents()
		window.SwapBuffers()

	}

	glfw.Terminate()

}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

// initOpenGL initializes OpenGL
func initOpenGL() {
	var nrAttributes int32
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))

	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	log.Println("nrAttributes", nrAttributes)
	log.Println("OpenGL version", version)
}

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

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
