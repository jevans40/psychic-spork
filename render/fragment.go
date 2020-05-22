package render

var fragmentShader = `

#version 450 core

//layout(location = 0) uniform sampler2D texture1;
//uniform sampler2D texture2;

in vec2 TexPos;
in vec4 Color;
in flat int TexMap;

out vec4 fragcolor;

void main(){
  fragcolor = Color;
}
` + "\x00"
