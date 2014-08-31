#version 440

layout(points) in;
layout(triangle_strip, max_vertices=256) out;

uniform mat4 projection;

smooth out float intensivity;

const float SIDES = 32;
const float PI = 3.1415926535897932384626433832795;

void main() {
	vec4 inPosition = gl_in[0].gl_Position;
	float pointSize = gl_in[0].gl_PointSize;
	for (int i = 0; i < SIDES; i++) {
		float theta = PI * 2 / SIDES * i;
		vec4 off = vec4(cos(theta), -sin(theta), 0, 0) * pointSize * projection;
		intensivity = 0;
		gl_Position = inPosition + off;
		EmitVertex();

		intensivity = 1;
		gl_Position = inPosition;
		EmitVertex();
	}
	intensivity = 0;
	gl_Position = inPosition + vec4(1, 0, 0, 0) * pointSize * projection;
	EmitVertex();

	EndPrimitive();
}
