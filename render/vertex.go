package render

var vertexShader = `
#version 330 core
layout(location = 0) in vec3 vertexPosition_modelspace;
out vec3 fragmentColor;
void main(){
    gl_Position.xyz = vertexPosition_modelspace;
    gl_Position.w = 1.0;
    fragmentColor = vertexPosition_modelspace;
  }
` + "\x00"
