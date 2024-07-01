#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

// texture samplers
uniform sampler2D texture1;
uniform bool Water;

void main()
{
  if (Water) {
    FragColor = vec4(0.0, 0.0, 1.0, 0.8);
  } else {
    FragColor = texture(texture1, TexCoord);
  }
}
