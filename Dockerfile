FROM golang:1.17-alpine
COPY ./vertexai-imagen vertexai-imagen
ENTRYPOINT ["./vertexai-imagen"]