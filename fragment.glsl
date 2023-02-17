#version 330 core

in vec2 texCoords;

uniform vec4 uColor;
uniform vec4 uOffset;
uniform sampler2D uImage;
uniform int uType;

out vec4 color;

void main () {
  if (uType == 0){
    color = uColor;    
  }
  else {
    vec2 t = texCoords;
    color = texture (uImage, t * uOffset.zw + uOffset.xy) * uColor;    
  }
}
