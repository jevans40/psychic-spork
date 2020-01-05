package graphics

var fragmentShader = `
#version 330 core
out vec3 color;
void main(){
  color = vec3(1,0,0);
}
` + "\x00"
