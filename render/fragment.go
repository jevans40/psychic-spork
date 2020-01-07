package render

var fragmentShader = `
#version 330 core
out vec3 color;
in vec3 fragmentColor;
void main(){
  color = fragmentColor;
}
` + "\x00"
