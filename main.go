package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"helloOpenGLWindow/shader"
	"log"
	"runtime"
)

var triangle1 = []float32{
	//first triangle
	//positions
	0, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
}

var triangle2 = []float32{
	//first triangle
	//positions       //colors
	0.5,  -0.5, 0.0,  1.0, 0.0, 0.0,
	-0.5, -0.5, 0.0,  0.0, 1.0, 0.0,
	0.0,  0.5,  0.0,  0.0, 0.0, 1.0,
}

var triangle3 = []float32{
	//third triangle
	1, 1, 0,
	.5, .5, 0,
	.5, 0, 0,
}

var triangle4 = []float32{
	//first triangle
	//positions      //colors
	1, 1, 0, 1.0, 0.0, 0.0,
	.5, .5, 0, 0.0, 1.0, 0.0,
	1, 0, 0, 0.0, 0.0, 1.0,
}

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
)

func main() {

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

	log.Println("Creating and Compiling Shaders")
	vertexShaderPath := "./shaders/shader.vs"
	fragmentShaderPath := "./shaders/shader.fs"
	shader := shader.New(vertexShaderPath, fragmentShaderPath)

	//vao1 := makeVao(triangle1)
	vao2 := makeVao(triangle2)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, shader.VertexShaderCompiled)
	gl.AttachShader(prog, shader.FragmentShaderCompiled)
	gl.LinkProgram(prog)

	prog2 := gl.CreateProgram()
	gl.AttachShader(prog2, shader.VertexShaderCompiled)
	gl.AttachShader(prog2, shader.FragmentShaderCompiled)
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

		/*
		gl.UseProgram(prog)
		timeValue := glfw.GetTime()
		greenValue := (math.Sin(timeValue) / 2) + 0.5
		vertexColorLocation := gl.GetUniformLocation(prog, gl.Str("ourColor\x00"))
		gl.Uniform4f(vertexColorLocation, 0.0, float32(greenValue), 0.0, 1.0)
		//gl.BindVertexArray(vao1)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		*/

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

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	var vao uint32
	var stride int32

	//points only 9
	//points and colors 18
	stride = int32(4 * len(points) / 3)
	println("stride: ", stride)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))
	println("triangle length: ", len(points))
	if len(points) >= 18 {
		log.Println("In if")
		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*4))
	}
	return vao
}
