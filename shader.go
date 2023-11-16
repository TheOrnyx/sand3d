package main

import (
	"fmt"
	"strings"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type shader struct {
	ID uint32 //program ID
}

func NewShader(vertexShaderSource, fragmentShaderSource []byte, name string) (*shader, error) {
	newShader := new(shader)
	
	vertexShader, err := compileShader(string(vertexShaderSource), gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile vertex Shader at %v\n%v", name, err)
	}
	
	fragmentShader, err := compileShader(string(fragmentShaderSource), gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile fragment Shader at %v\n%v", name, err)
	}
	
	newShader.ID = gl.CreateProgram()
	gl.AttachShader(newShader.ID, vertexShader)
	gl.AttachShader(newShader.ID, fragmentShader)
	gl.LinkProgram(newShader.ID)

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	
	return newShader, nil
}

func compileShader(shaderCode string, shaderType uint32) (uint32, error) {
	var shader uint32

	shader = gl.CreateShader(shaderType)
	cSources, free := gl.Strs(shaderCode+"\x00")
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile\n%v\n%v", shaderCode, log)
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
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), intVal)
}

func (s *shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}

func (s *shader) SetVec2(name string, value *mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *shader) SetVec2f(name string, x, y float32) {
	gl.Uniform2f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y)
}

func (s *shader) SetVec3(name string, value *mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *shader) SetVec3f(name string, x, y, z float32) {
	gl.Uniform3f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y, z)
}

func (s *shader) SetVec4(name string, value *mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *shader) SetVec4f(name string, x, y, z, w float32) {
	gl.Uniform4f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y, z, w)
}

func (s *shader) SetMat2(name string, mat *mgl32.Mat2) {
	gl.UniformMatrix2fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &(*mat)[0])
}

func (s *shader) SetMat3(name string, mat *mgl32.Mat3) {
	gl.UniformMatrix3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &(*mat)[0])
}

func (s *shader) SetMat4(name string, mat *mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &(*mat)[0])
}
