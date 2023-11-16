#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

// texture samplers
uniform sampler2D texture1;
uniform bool white;

void main()
{
  if (white) {
    FragColor = vec4(1.0, 1.0, 1.0, 1.0);
  } else {
    FragColor = texture(texture1, TexCoord);
  }
}
