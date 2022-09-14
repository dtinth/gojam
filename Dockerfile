FROM golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apt-get update && apt-get install -y automake
RUN git clone https://github.com/jamulussoftware/jamulus.git /opt/jamulus
RUN bash -c 'cd /opt/jamulus/libs/opus && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'

COPY . .
ENV CGO_CFLAGS=-I/opt/jamulus/libs/opus/include
ENV CGO_LDFLAGS=-L/opt/jamulus/libs/opus/.libs
RUN go build -v ./cmd/gojamclient
