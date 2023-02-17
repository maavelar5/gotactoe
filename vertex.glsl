#version 330 core

layout (location = 0) in vec2 aPos;

uniform mat4 uProjection;
uniform mat4 uModel;

out vec2 texCoords;

void main () {
  texCoords = aPos;
  gl_Position = uProjection * uModel * vec4(aPos, 0.0f, 1.0f);
}
