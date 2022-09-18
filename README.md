# gojam

Go implementation of a Jamulus client.

To be able to use this client, we need to build a custom version of Opus from Jamulus source tree.

```sh
git clone https://github.com/jamulussoftware/jamulus.git
bash -c 'cd jamulus/libs/opus && ./configure --enable-static --disable-shared --enable-custom-modes --disable-hardening && make'
```

Then provide these environment variables to `go build` or `go run` in order to successfully build and run your code:

```sh
export CGO_CFLAGS="-I$PWD/jamulus/libs/opus/include"
export CGO_LDFLAGS="-L$PWD/jamulus/libs/opus/.libs"
```

## Usage

A simple CLI, `gojamclient` is provided that can connect to a Jamulus server and send received audio to another TCP server as raw PCM samples (16-bit, 48kHz, 2 channels).

First, you need to run a TCP server that will receive the audio samples. For example, you can use a combination of `nc` and `play` (from `sox` package) to listen to the audio:

```sh
nc -l 22222 | play -r 48000 -c 2 -b 16 -e signed -t raw -
```

Then you can run `gojamclient`:

```sh
go run ./cmd/gojamclient/ -server <ip>:<port> -pcmout localhost:22222
```

### REST API

When you run an API server, e.g. `-apiserver localhost:9999`, you can control the client, for example, using [HTTPie](https://httpie.io/):

```sh
# Get the client info
http get localhost:9999/channel-info

# Update the client info
http patch localhost:9999/channel-info name=newname
```
