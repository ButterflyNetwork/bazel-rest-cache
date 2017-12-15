FROM golang:latest

EXPOSE 8080
ENV REDIS_ADDR $REDIS_ADDR

RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app

RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT go-wrapper run --port=8080 --redis_addr="$REDIS_ADDR"
