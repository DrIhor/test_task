FROM golang:latest

# create work dir
RUN mkdir -p /go/src/github.com/DrIhor/test_task/
ADD .  /go/src/github.com/DrIhor/test_task/
WORKDIR  /go/src/github.com/DrIhor/test_task/

# build program
RUN cp -a ./config/ ./bin/
RUN go build -o ./bin/main/ ./cmd/http

# start program
WORKDIR /go/src/github.com/DrIhor/test_task/bin/
CMD ["./main"]