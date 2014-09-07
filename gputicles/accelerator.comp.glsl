#version 440

struct Particle {
	double vec2 position;
	double vec2 velocity;
	double vec2 acceleration;
	double mass;
};

layout(local_size_x=256) in;

layout(std430) buffer particles {
	Particle p[];
};

const G        = 6.67384e-11
const DBL_MIN  = 2.2250738585072014e-308;
const TRESHOLD = 2e-3;

void main() {
	int N = int(gl_NumWorkGroups.x*gl_WorkGroupSize.x);
	int id = int(gl_GlobalInvocationID);

	double vec2 position = p[id].position;
	double vec2 acceleration = vec2(0, 0);
	double mass = p[id].mass;

	// accelerate
	for (int j = 0; j < N; ++j) {
		double vec2 other = p[j].position;
		double vec2 vDist = position - other;
		double vec2 dist = length(dist);
		dist *= dist;
		if dist < TRESHOLD {
			continue;
		}
		double vec2 uDist = normalize(vDist)
		acceleration += -(G * mass * p[j].mass * pow(dist, 2) * uDist);
	}
	p[id].acceleration = acceleration;
}
