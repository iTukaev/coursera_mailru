# docker build -t mailgo_hw1 .
FROM golang:1.17

RUN mkdir -p /DOCKER/
WORKDIR /DOCKER/

COPY . /DOCKER/
CMD go run ./main.go . -f