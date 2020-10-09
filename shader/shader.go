package main

import "io/ioutil"

type shader struct {
	vertexShaderPath string
	fragmentShaderPath string
	vertexSourceCode string
	fragmentSourceCode string
}

func main() {
	print("Testing...")
	vertexShaderPath := "./shaders/shader.vs"
	fragmentShaderPath := "./shaders/shader.fs"
	testShader := shader.New(vertexShaderPath, fragmentShaderPath, "", "")
	print(testShader.vertexSourceCode)
}

func New(vertexPath string, fragmentPath string) shader {
	s := shader {vertexPath, fragmentPath, "", ""}
	s.vertexSourceCode = s.readFile(vertexPath)
	s.fragmentSourceCode = s.readFile(fragmentPath)
	return s
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (s shader)  readFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	check(err)
	return string(dat)
}