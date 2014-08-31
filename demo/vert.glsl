#version 440

uniform mat4 projection;

in vec2 position;

void main() {
	gl_Position = vec4(position, 0, 1) * projection;
	gl_PointSize = 0.004;
}
