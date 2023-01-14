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

### EventStream (requires `-apiserver`)

When you run an API server, you can also subscribe to [server-sent events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) to receive events from the server.

```sh
# Subscribe to events
curl localhost:9999/events
```

Upon subscribing, you will receive objects with these keys (all of them are optional):

- `clients` - Emitted when the list of clients changes. The value is an array of client objects.
- `levels` - Emitted when the list of client levels changes. The value is an array of numbers.
- `newChatMessage` - Emitted when a new chat message is received. The value is a chat message object (same as the one returned by `/chat` endpoint).

### MP3 stream (requires `-apiserver` and `ffmpeg`)

When you run an API server, you can also stream the audio as MP3.

```sh
# Play the MP3 stream
ffplay -i http://localhost:9999/mp3
```

### Listener client

Run `gojam` with `-apiserver` on port 9999 and with `-mp3` flag:

```sh
go run ./cmd/gojamclient/ -apiserver :9999 -mp3 -server 127.0.0.1:22124
```

You can then open `visualizer.html` in your browser to listen to the audio stream being sent to it.

### Streaming Jamulus audio to Discord

There is a separate project [pcm2discord](https://github.com/dtinth/pcm2discord) that receives raw PCM samples via TCP, and sends them to Discord. A Python script [jamulus-discord-glue](https://github.com/dtinth/jamulus-discord-glue) is used to [glue](https://en.wikipedia.org/wiki/Glue_code) these 2 systems together.
