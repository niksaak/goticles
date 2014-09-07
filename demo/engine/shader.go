package engine

import (
	"fmt"
	"github.com/go-gl/glow/gl-core/4.4/gl"
	"io"
	"os"
	"strings"
)

type ShaderKind uint32

const (
	VertexShader         ShaderKind = gl.VERTEX_SHADER
	TessControlShader               = gl.TESS_CONTROL_SHADER
	TessEvaluationShader            = gl.TESS_EVALUATION_SHADER
	GeometryShader                  = gl.GEOMETRY_SHADER
	FragmentShader                  = gl.FRAGMENT_SHADER
	ComputeShader                   = gl.COMPUTE_SHADER
)

type ShaderProgram struct {
	shaders  map[ShaderKind][]uint32
}

func NewShaderProgram() *ShaderProgram {
	p := new(ShaderProgram)
	p.shaders = make(map[ShaderKind][]uint32)
	return p
}

func (p *ShaderProgram) AddShader(source string, kind ShaderKind) error {
	shader := gl.CreateShader(uint32(kind))

	csource := gl.Str(source + "\x00")
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return fmt.Errorf("failed to compile %v shader:\n%s",
			kindString(kind), log)
	}
	p.shaders[kind] = append(p.shaders[kind], shader)
	return nil
}

const _READ_SHADER_BUFFER_LEN = 1024

func (p *ShaderProgram) ReadShader(r io.Reader, kind ShaderKind) error {
	var buffer [_READ_SHADER_BUFFER_LEN]byte
	var source string
	for {
		sourceLen, err := r.Read(buffer[:])
		if err != nil {
			return err
		}
		source += string(buffer[:sourceLen])
		if sourceLen != _READ_SHADER_BUFFER_LEN {
			break
		}
	}
	return p.AddShader(source, kind)
}

func (p *ShaderProgram) ReadShaderFile(path string, kind ShaderKind) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return p.ReadShader(file, kind)
}

func (p *ShaderProgram) Link() (uint32, error) {
	defer func(){
		for _, shaders := range p.shaders {
			for _, shader := range shaders {
				gl.DeleteShader(shader)
			}
		}
		*p = ShaderProgram{}
	}()

	program := gl.CreateProgram()
	for _, shaders := range p.shaders {
		for _, shader := range shaders {
			gl.AttachShader(program, shader)
		}
	}
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program:\n%s", log)
	}
	return program, nil
}

func kindString(kind ShaderKind) string {
	switch kind {
	case gl.VERTEX_SHADER:
		return "vertex"
	case gl.TESS_CONTROL_SHADER:
		return "tesselation control"
	case gl.TESS_EVALUATION_SHADER:
		return "tesselation evaluation"
	case gl.GEOMETRY_SHADER:
		return "geometry"
	case gl.FRAGMENT_SHADER:
		return "fragment"
	case gl.COMPUTE_SHADER:
		return "compute"
	default:
		return "unknown"
	}
}
