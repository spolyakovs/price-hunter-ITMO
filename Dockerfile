FROM golang:1.17
WORKDIR /go/src/github.com/redis_docker
COPY . .
RUN apt update
RUN apt install make
RUN go get -d ./...
CMD make start CONFIG_PATH=configs/docker.toml
EXPOSE 8000
