#version 440

smooth in float intensivity;

void main() {
	gl_FragColor = vec4(0.3, 0.9, 0.3, 0.2 * intensivity);
}
