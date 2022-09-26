# gojam

An implementation of a Jamulus client written in Go, intended for headless streaming use-cases.

## Building a Docker image

```sh
docker build -t gojam .
```

## Building locally

Building this project requires custom modes build of Opus from Jamulus source tree.

```sh
git clone https://github.com/jamulussoftware/jamulus.git
bash -c 'cd jamulus/libs/opus && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'
```

Provide these environment variables to `go build` or `go run` in order to successfully build and run this project:

```sh
export CGO_CFLAGS="-I$PWD/jamulus/libs/opus/include"
export CGO_LDFLAGS="-L$PWD/jamulus/libs/opus/.libs"
```

## Usage

A simple CLI, `gojamclient` is provided. It can connect to a Jamulus server and send received audio to another TCP server as raw PCM samples (16-bit signed, little endian, 48kHz, 2 channels).

First, you need to run a TCP server that will receive the audio samples. For example, you can use a combination of `nc` and `play` (from `sox` package) to listen to the audio:

```sh
nc -l 22222 | play -r 48000 -c 2 -b 16 -e signed -t raw -
```

Then you can run `gojamclient`:

```sh
go run ./cmd/gojamclient/ -server <ip>:<port> -pcmout localhost:22222
```

`gojamclient` comes with an extra large network jitter buffer, so you can listen to Jamulus servers on your Wi-Fi connection.

### REST API

When you run an API server, e.g. `-apiserver localhost:9999`, you can control the client, for example, using [HTTPie](https://httpie.io/):

```sh
# Get the client info
http get localhost:9999/channel-info

# Update the client info
http patch localhost:9999/channel-info name=newname

# Retrieve latest 100 chat messages
http get localhost:9999/chat

# Send a chat message
http post localhost:9999/chat message="hello world"
```

### Streaming Jamulus audio to Discord

There is a separate project [pcm2discord](https://github.com/dtinth/pcm2discord) that receives raw PCM samples via TCP, and sends them to Discord. A Python script [jamulus-discord-glue](https://github.com/dtinth/jamulus-discord-glue) is used to [glue](https://en.wikipedia.org/wiki/Glue_code) these 2 systems together.
