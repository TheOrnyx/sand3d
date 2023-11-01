package main

import (
	_"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type shader struct {
	ID uint32 //program ID
}

func NewShader(vertexPath, fragmentPath string) *shader {
	var vertexCode, fragmentCode string
	newShader := new(shader)
	
	vertexShaderSource, err := os.ReadFile(vertexPath)
	handleError(err)
	fragmentShaderSource, err := os.ReadFile(fragmentPath)
	handleError(err)

	vertexCode = string(vertexShaderSource)
	fragmentCode = string(fragmentShaderSource)
	
	vertexShader, err := compileShaders(vertexCode, gl.VERTEX_SHADER)
	handleError(err)
	fragmentShader, err := compileShaders(fragmentCode, gl.FRAGMENT_SHADER)
	handleError(err)


	newShader.ID = gl.CreateProgram()
	gl.AttachShader(newShader.ID, vertexShader)
	gl.AttachShader(newShader.ID, fragmentShader)
	gl.LinkProgram(newShader.ID)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	
	return newShader
}

func compileShaders(shaderCode string, shaderType uint32) (uint32, error) {
	var shader uint32

	shader = gl.CreateShader(shaderType)
	cSources, free := gl.Strs(shaderCode)
	gl.ShaderSource(shader, 1, cSources, nil)
	gl.CompileShader(shader)
	free()

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", shaderCode, log)
	}
	
	return shader, nil
}

func (s *shader) use() {
	gl.UseProgram(s.ID)
}

func (s *shader) SetBool(name string, value bool) {
	var intVal int32
	if value {
		intVal = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name)), intVal)
}

func (s *shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name)), value)
}

func (s *shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name)), value)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
