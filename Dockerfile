FROM golang:1.19 AS opus
WORKDIR /usr/src/app
RUN apt-get update && apt-get install -y automake libtool
RUN git clone https://github.com/jamulussoftware/jamulus.git /opt/jamulus
RUN bash -c 'cd /opt/jamulus/libs/opus && autoreconf -f -i && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'

FROM golang:1.19
RUN if [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then \
  wget https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0.1/linux-arm64 -O /usr/local/bin/ffmpeg; \
  else \
  wget https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0.1/linux-x64 -O /usr/local/bin/ffmpeg; \
  fi && chmod +x /usr/local/bin/ffmpeg
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY --from=opus /opt/jamulus/ /opt/jamulus/
COPY . .
ENV CGO_CFLAGS=-I/opt/jamulus/libs/opus/include
ENV CGO_LDFLAGS=-L/opt/jamulus/libs/opus/.libs
RUN go build -v ./cmd/gojamclient
RUN go install -v ./cmd/gojamclient
