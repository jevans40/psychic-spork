package render

var fragmentShader = `

#version 450 core

uniform sampler2D texture1;
//uniform sampler2D texture2;

in vec2 TexPos;
in vec4 Color;
in flat int TexMap;

out vec4 fragcolor;

void main(){
  if (TexMap > 0){
    fragcolor = texture(texture1,TexPos) * Color;
  }
  else {
    fragcolor = Color;
  }
}
` + "\x00"
